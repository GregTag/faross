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

func (r *RepositoryStore) GetUnquarantined(instanceId, name string, since time.Time) (*[]entity.UnquarantineEntry, error) {
	var entries []entity.UnquarantineEntry
	err := r.db.Model(&entity.Repository{InstanceID: instanceId, RepositoryDTO: entity.RepositoryDTO{Name: name}}).
		Association("UnquarantineList").
		Find(&entries, "created_at >= ?", since)
	return &entries, err
}

func (r *RepositoryStore) AppendPackages(instanceId, name string, packages []*entity.Package) error {
	return r.db.
		Model(&entity.Repository{}).
		Where("instance_id = ?", instanceId).
		Where("name = ?", name).
		Association("Packages").
		Append(packages)
}
