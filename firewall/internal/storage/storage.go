package storage

import "gorm.io/gorm"

type Storage struct {
	Repository RepositoryStore
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Repository: NewRepositoryStore(db),
	}
}
