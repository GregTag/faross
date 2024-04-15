package storage

import (
	"firewall/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryStore struct {
	db *gorm.DB
}

func NewRepositoryStore(db *gorm.DB) RepositoryStore {
	return RepositoryStore{db}
}

func (r *RepositoryStore) Save(entries any) error {
	return r.db.Model(&entity.Repository{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "instance_id"}, {Name: "name"}},
		UpdateAll: true,
	}).Create(entries).Error
}

func (r *RepositoryStore) GetByInstanceId(instanceId string) (*[]entity.Repository, error) {
	return r.GetByInstanceIdAndSince(instanceId, time.Time{})
}

func (r *RepositoryStore) GetByInstanceIdAndSince(instanceId string, since time.Time) (*[]entity.Repository, error) {
	var entries []entity.Repository
	err := r.db.Model(&entity.Repository{}).Where("instance_id = ?", instanceId).Where("updated_at >= ?", since).Scan(&entries).Error
	return &entries, err
}

func (r *RepositoryStore) GetByInstanceAndName(instance, name string) (*entity.Repository, error) {
	var entry entity.Repository
	err := r.db.Where("instance_id = ?", instance).Where("name = ?", name).Take(&entry).Error
	return &entry, err
}

func (r *RepositoryStore) GetById(id uint) (*entity.Repository, error) {
	entry := entity.Repository{Model: gorm.Model{ID: id}}
	err := r.db.Take(&entry).Error
	return &entry, err
}

func (r *RepositoryStore) Load(instanceId, name string) (*entity.Repository, error) {
	var entry entity.Repository
	err := r.db.Where("instance_id = ?", instanceId).Where("name = ?", name).Preload("Packages").Take(&entry).Error
	return &entry, err
}

func (r *RepositoryStore) GetUnquarantined(id uint, since time.Time) ([]*entity.Package, error) {
	var entries []*entity.Package
	err := r.db.Model(&entity.Repository{Model: gorm.Model{ID: id}}).
		Where("updated_at >= ?", since).
		Where("state = ?", entity.Unquarantined).
		Association("Packages").
		Find(&entries)
	return entries, err
}

func (r *RepositoryStore) AppendPackages(id uint, packages []*entity.Package) error {
	return r.db.Model(&entity.Repository{Model: gorm.Model{ID: id}}).
		Association("Packages").
		Append(&packages)
}
