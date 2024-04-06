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

func (s *Service) GetConfiguredRepositories(instance, sinceStr string) (*[]entity.RepositoryDTO, error) {
	sinceTimestamp, err := strconv.ParseInt(sinceStr, 10, 64)
	if err != nil {
		return nil, err
	}
	since := time.Unix(sinceTimestamp, 0)

	entries, err := s.storage.Repository.GetByInstanceIdAndSince(instance, since)
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
