package service

import (
	"encoding/json"
	"errors"
	"faross/gatherlaunch"
	"firewall/internal/entity"
	"firewall/pkg/config"
	"firewall/pkg/parsers"
	"log"
	"sync"
	"time"

	"github.com/package-url/packageurl-go"
)

type EvalDataRequest struct {
	Format   string `json:"format"`
	Pathname string `json:"pathname"`
	Hash     string `json:"hash"`
}

func (s *Service) runGatherLaunch(purl *packageurl.PackageURL) (*entity.Package, error) {
	resp, err := gatherlaunch.Scan(*purl)
	if err != nil {
		return nil, err
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	respStr := string(respBytes)

	log.Println("Scanned response for ", purl, ": ", respStr)

	var state entity.State
	if resp.IsQuarantined {
		state = entity.Quarantined
	} else {
		state = entity.Healthy
	}

	result := entity.Package{
		Purl:       purl.ToString(),
		State:      state,
		FinalScore: resp.Score,
		Report:     respStr,
	}

	return &result, nil
}

var maxPendingDurating time.Duration

func (s *Service) evaluate(purl *packageurl.PackageURL, pathname string) (*entity.Package, error) {
	if maxPendingDurating == 0 {
		maxPendingDurating = config.Koanf.Duration("maxPendingDuration")
	}

	pkg, err := s.storage.Package.GetOrInsertPending(purl.ToString(), pathname)
	duration := time.Since(pkg.CreatedAt)
	if err == nil && (pkg.CreatedAt != pkg.UpdatedAt) && (pkg.State != entity.Pending || duration < maxPendingDurating) {
		if pkg.State == entity.Pending {
			log.Printf("Package %s is pending with duration %s\n", pkg.Purl, duration)
			return nil, entity.ErrPending
		}
		return pkg, nil
	}

	pkg, err = s.runGatherLaunch(purl)
	if err != nil {
		return nil, err
	}
	pkg.Pathname = pathname

	err = s.storage.Package.Save(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func (s *Service) EvaluatePurl(purlStr string) (*entity.Package, error) {
	purl, err := packageurl.FromString(purlStr)
	if err != nil {
		return nil, err
	}

	return s.evaluate(&purl, "")
}

func (s *Service) EvaluatePathname(format, pathname string) (*entity.Package, error) {
	purl, err := parsers.PathnameToPurl(format, pathname)
	if err != nil {
		return nil, err
	}

	return s.evaluate(purl, pathname)
}

type runResult struct {
	index int
	pkg   *entity.Package
}

func (s *Service) runEvals(components []EvalDataRequest) ([]runResult, error) {
	results := make(chan runResult, len(components))
	errs := make(chan error, len(components))

	var wg sync.WaitGroup
	wg.Add(len(components))
	for index, component := range components {
		go func(index int, component EvalDataRequest) {
			defer wg.Done()
			pkg, err := s.EvaluatePathname(component.Format, component.Pathname)
			if err != nil {
				errs <- err
				return
			}
			results <- runResult{
				index: index,
				pkg:   pkg,
			}
		}(index, component)
	}
	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	runResults := make([]runResult, 0, len(components))
	runErrors := make([]error, 0, len(components))

	for i := 0; i < len(components); i++ {
		select {
		case result, ok := <-results:
			if ok {
				runResults = append(runResults, result)
			}
		case err, ok := <-errs:
			if ok {
				runErrors = append(runErrors, err)
			}
		}
	}
	return runResults, errors.Join(runErrors...)
}

func (s *Service) EvalRequest(instance, name string, components []EvalDataRequest) ([]map[string]any, error) {
	repos, err := s.storage.Repository.GetByInstanceAndName(instance, name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	runResults, err := s.runEvals(components)
	if err != nil {
		log.Println("Errors in evaluation: ", err)
		// No break
	}

	var evalResults []map[string]any
	var pkgs = make([]entity.Package, len(components))
	for _, res := range runResults {
		evalResults = append(evalResults,
			map[string]any{
				"requestIndex": res.index,
				"quarantine":   res.pkg.State == entity.Quarantined,
			})
		pkgs[res.index] = *res.pkg
	}

	err = s.storage.Repository.AppendPackages(repos.ID, pkgs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return evalResults, nil
}

func (s *Service) Unquarantine(purl, comment string) error {
	err := s.storage.Package.Unquarantine(purl, comment)
	if err != nil {
		log.Println("Error in unqarantine: ", err)
	}
	return err
}

func (s *Service) GetPackage(purl string) (*entity.Package, error) {
	return s.storage.Package.GetByPurl(purl)
}

func preparePackages(pkgs []entity.Package) []map[string]any {
	prepared := make([]map[string]any, 0, len(pkgs))
	for _, pkg := range pkgs {
		prepared = append(prepared, map[string]any{
			"purl":       pkg.Purl,
			"state":      pkg.State.ToSring(),
			"score":      pkg.FinalScore,
			"comment":    pkg.Comment,
			"changed_at": pkg.UpdatedAt,
		})
	}
	return prepared
}

func (s *Service) GetAll() ([]map[string]any, error) {
	pkgs, err := s.storage.Package.GetAll()
	if err != nil {
		return nil, err
	}
	return preparePackages(pkgs), nil
}

func (s *Service) ChangeComment(purl, comment string) error {
	return s.storage.Package.UpdateComment(purl, comment)
}
