package storage

import (
	"errors"
	"firewall/internal/entity"

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

func (s *PackageStore) GetByPurl(purl string) (*entity.Package, error) {
	var entry entity.Package
	err := s.db.Model(&entity.Package{}).Where("purl = ?", purl).Take(&entry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrPackageNotFound
	}
	return &entry, err
}

func (s *PackageStore) GetOrInsertPending(purl, pathname string) (*entity.Package, error) {
	var entry entity.Package
	result := s.db.
		Where(entity.Package{Purl: purl}).
		Assign(entity.Package{Pathname: pathname}).
		Attrs(entity.Package{State: entity.Pending}).
		FirstOrCreate(&entry)
	return &entry, result.Error
}

func (s *PackageStore) Unquarantine(purl, comment string) error {
	result := s.db.Model(&entity.Package{}).
		Where("purl = ?", purl).
		Where("state = ?", entity.Quarantined).
		Updates(entity.Package{State: entity.Unquarantined, Comment: comment})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrNothingUnquarantined
	}
	return nil
}

func (s *PackageStore) GetAll() ([]entity.Package, error) {
	var entries []entity.Package
	err := s.db.Find(&entries).Error
	return entries, err
}

func (s *PackageStore) UpdateComment(purl, comment string) error {
	result := s.db.Model(&entity.Package{}).Where("purl = ?", purl).Update("comment", comment)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrPackageNotFound
	}
	return nil
}
