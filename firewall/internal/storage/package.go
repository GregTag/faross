package storage

import (
	"errors"
	"firewall/internal/entity"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PackageStore struct {
	db *gorm.DB
}

func NewPackageStore(db *gorm.DB) PackageStore {
	return PackageStore{db: db}
}

func (s *PackageStore) Save(pkg *entity.Package) error {
	return s.db.Model(&entity.Package{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "purl"}},
		UpdateAll: true,
	}).Create(pkg).Error
}

func (s *PackageStore) TryGetByPurl(purl string) (*entity.Package, error) {
	var entry entity.Package
	err := s.db.Model(&entity.Package{}).Where("purl = ?", purl).Take(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entry, err
}

func (s *PackageStore) Unquarantine(purl string) error {
	result := s.db.Model(&entity.Package{}).
		Where("purl = ?", purl).
		Where("state = ?", entity.Quarantined).
		Update("state", entity.Unquarantined)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no package unquarantined")
	}
	return nil
}

func (s *PackageStore) GetAll() ([]entity.Package, error) {
	var entries []entity.Package
	err := s.db.Find(&entries).Error
	return entries, err
}
