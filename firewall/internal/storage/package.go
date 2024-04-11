package storage

import (
	"errors"
	"firewall/internal/entity"

	"gorm.io/gorm"
)

type PackageStore struct {
	db *gorm.DB
}

func NewPackageStore(db *gorm.DB) PackageStore {
	return PackageStore{db: db}
}

func (s *PackageStore) Save(pkg *entity.Package) error {
	return s.db.Save(pkg).Error
}

func (s *PackageStore) TryGetByPurl(purl string) (*entity.Package, error) {
	var entry entity.Package
	err := s.db.Model(&entity.Package{}).Where("purl = ?", purl).Take(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entry, err
}

func (s *PackageStore) Unquarantine(pkg *entity.Package) error {
	pkg.State = entity.Unquarantined
	return s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(pkg).Error
		// TODO: insert UnquarantineEntry to each repository for this package
	})
}
