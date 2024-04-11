package storage

import "gorm.io/gorm"

type Storage struct {
	Repository RepositoryStore
	Package    PackageStore
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Repository: NewRepositoryStore(db),
		Package:    NewPackageStore(db),
	}
}
