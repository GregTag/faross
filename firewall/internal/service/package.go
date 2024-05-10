package service

import (
	"encoding/json"
	"errors"
	"faross/gatherlaunch"
	"firewall/internal/entity"
	"firewall/pkg/parsers"
	"log"
	"sync"

	"github.com/package-url/packageurl-go"
)

type EvalDataRequest struct {
	Format   string `json:"format"`
	Pathname string `json:"pathname"`
	Hash     string `json:"hash"`
}

func (s *Service) requestEvaluatePurl(purl *packageurl.PackageURL) (*entity.Package, error) {
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

	// TODO: get decision from report
	state := entity.Quarantined
	// switch resp.Decision {
	// case "quarantine":
	// 	state = entity.Quarantined
	// case "healthy":
	// 	state = entity.Healthy
	// default:
	// 	err = fmt.Errorf("invalid decision")
	// 	log.Println(err)
	// 	return nil, err
	// }

	final_score, ok := resp["final_score"].(float64)
	if !ok {
		return nil, errors.New("no final score in report")
	}

	result := entity.Package{
		Purl:       purl.ToString(),
		State:      state,
		FinalScore: final_score,
		Report:     respStr,
	}

	return &result, nil
}

func (s *Service) EvaluatePurl(purlStr string) (*entity.Package, error) {
	purl, err := packageurl.FromString(purlStr)
	if err != nil {
		return nil, err
	}

	pkg, err := s.requestEvaluatePurl(&purl)
	if err != nil {
		return nil, err
	}

	err = s.storage.Package.Save(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func (s *Service) EvaluatePathname(format, pathname string) (*entity.Package, error) {
	purl, err := parsers.PathnameToPurl(format, pathname)
	if err != nil {
		return nil, err
	}

	pkg, err := s.requestEvaluatePurl(purl)
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

func (s *Service) Unquarantine(purl string) error {
	err := s.storage.Package.Unquarantine(purl)
	if err != nil {
		log.Println("Error in unqarantine: ", err)
	}
	return err
}

func (s *Service) GetPackage(purl string) (*entity.Package, error) {
	return s.storage.Package.TryGetByPurl(purl)
}

func (s *Service) GetAll() ([]map[string]any, error) {
	pkgs, err := s.storage.Package.GetAll()
	if err != nil {
		return nil, err
	}
	response := make([]map[string]any, 0, len(pkgs))
	for _, pkg := range pkgs {
		response = append(response, map[string]any{
			"purl":  pkg.Purl,
			"state": pkg.State,
		})
	}
	return response, nil
}
