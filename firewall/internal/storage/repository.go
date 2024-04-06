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
