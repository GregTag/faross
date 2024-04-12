package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"firewall/internal/entity"
	"firewall/pkg/config"
	"firewall/pkg/parsers"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type EvalDataRequest struct {
	Format   string `json:"format"`
	Pathname string `json:"pathname"`
	Hash     string `json:"hash"`
}

func (s *Service) requestEvaluatePurl(purl string) (*entity.Package, error) {
	pkg, err := s.storage.Package.TryGetByPurl(purl)
	if err != nil {
		return nil, err
	}
	if pkg != nil {
		return pkg, nil
	}

	scan_url := config.Koanf.String("gatherLaunch")

	resp, err := http.Post(scan_url, "text/plain", bytes.NewBufferString(purl))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Scanned response for ", purl, ": ", string(body))

	var response struct {
		Decision   string  `json:"decision"`
		FinalScore float32 `json:"final_score"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var state entity.State
	switch response.Decision {
	case "quarantine":
		state = entity.Quarantined
	case "healthy":
		state = entity.Healthy
	default:
		err = fmt.Errorf("invalid decision")
		log.Println(err)
		return nil, err
	}

	result := entity.Package{
		Purl:       purl,
		State:      state,
		FinalScore: response.FinalScore,
		Report:     string(body),
	}

	return &result, nil
}

func (s *Service) EvaluatePurl(purl string) (*entity.Package, error) {
	pkg, err := s.requestEvaluatePurl(purl)
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
	purlStruct, err := parsers.PathnameToPurl(format, pathname)
	if err != nil {
		return nil, err
	}
	purl := purlStruct.ToString()

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
	var pkgs = make([]*entity.Package, len(components))
	for _, res := range runResults {
		evalResults = append(evalResults,
			map[string]any{
				"requestIndex": res.index,
				"quarantine":   res.pkg.State == entity.Quarantined,
			})
		pkgs[res.index] = res.pkg
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
