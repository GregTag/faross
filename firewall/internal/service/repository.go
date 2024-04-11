package service

import (
	"firewall/internal/entity"
	"log"
	"strconv"
	"time"
)

func (s *Service) ConfigureRepositories(instance string, repositories *[]entity.RepositoryDTO) error {
	entries := make([]entity.Repository, 0, len(*repositories))
	for _, data := range *repositories {
		entries = append(entries, entity.Repository{
			InstanceID:    instance,
			RepositoryDTO: data,
		})
	}

	err := s.storage.Repository.Save(&entries)
	if err != nil {
		log.Printf("Error in configure repositories: %s\n", err)
		return err
	}

	return nil
}

func parseTimestamp(timestamp string) (*time.Time, error) {
	sinceTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	time := time.Unix(sinceTimestamp, 0)
	return &time, nil
}

func (s *Service) GetConfiguredRepositories(instance, sinceStr string) (*[]entity.RepositoryDTO, error) {
	since, err := parseTimestamp(sinceStr)
	if err != nil {
		return nil, err
	}
	entries, err := s.storage.Repository.GetByInstanceIdAndSince(instance, *since)
	if err != nil {
		log.Printf("Error in get configured repositories: %s\n", err)
		return nil, err
	}

	repos := make([]entity.RepositoryDTO, 0, len(*entries))
	for _, entry := range *entries {
		repos = append(repos, entry.RepositoryDTO)
	}
	return &repos, nil
}

func (s *Service) SetAuditEnable(instance, name string, enabled bool) (*entity.ApiRepository, error) {
	repos := entity.Repository{InstanceID: instance, RepositoryDTO: entity.RepositoryDTO{Name: name, AuditEnabled: enabled}}
	err := s.storage.Repository.Save(&repos)
	if err != nil {
		return nil, err
	}
	return repos.ToApiRepository(), nil
}

func (s *Service) SetQuarantineEnable(instance, name string, enabled bool) (*entity.ApiRepository, error) {
	repos := entity.Repository{InstanceID: instance, RepositoryDTO: entity.RepositoryDTO{Name: name, QuarantineEnabled: enabled}}
	err := s.storage.Repository.Save(&repos)
	if err != nil {
		return nil, err
	}
	return repos.ToApiRepository(), nil
}

func (s *Service) GetSummary(instance, name string) (*map[string]any, error) {
	repos, err := s.storage.Repository.Load(instance, name)
	if err != nil {
		return nil, err
	}

	quarantined := 0
	critical := 0
	severe := 0
	moderate := 0

	for _, pkg := range repos.Packages {
		if pkg.State == entity.Quarantined {
			quarantined += 1
		}
		level := pkg.FinalScore
		switch {
		case level >= 8:
			critical++
		case level >= 4:
			severe++
		case level >= 2:
			moderate++
		}
	}

	response := map[string]any{
		"affectedComponentCount":    critical + severe + moderate,
		"criticalComponentCount":    critical,
		"severeComponentCount":      severe,
		"moderateComponentCount":    moderate,
		"quarantinedComponentCount": quarantined,
		// "reportUrl" : "", // TODO
	}
	return &response, nil
}

func (s *Service) GetUnquarantined(instance, name, sinceStr string) (*[]string, error) {
	since, err := parseTimestamp(sinceStr)
	if err != nil {
		return nil, err
	}

	entries, err := s.storage.Repository.GetUnquarantined(instance, name, *since)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0, len(*entries))
	for _, entry := range *entries {
		list = append(list, entry.Pathname)
	}
	return &list, nil
}
