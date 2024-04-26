package entity

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	InstanceID string `gorm:"uniqueIndex:instance_name"`
	RepositoryDTO
	Packages []Package `gorm:"many2many:repository_package;"`
}

type RepositoryDTO struct {
	Name                                     string `json:"name" gorm:"uniqueIndex:instance_name"`
	Format                                   string `json:"format"`
	Type                                     string `json:"type"`
	AuditEnabled                             bool   `json:"auditEnabled"`
	QuarantineEnabled                        bool   `json:"quarantineEnabled"`
	PolicyCompliantComponentSelectionEnabled bool   `json:"policyCompliantComponentSelectionEnabled"`
	NamespaceConfusionProtectionEnabled      bool   `json:"namespaceConfusionProtectionEnabled"`
}

type ApiRepository struct {
	Name                                     string `json:"public_id"`
	Format                                   string `json:"format"`
	Type                                     string `json:"type"`
	AuditEnabled                             bool   `json:"auditEnabled"`
	QuarantineEnabled                        bool   `json:"quarantineEnabled"`
	PolicyCompliantComponentSelectionEnabled bool   `json:"policyCompliantComponentSelectionEnabled"`
	NamespaceConfusionProtectionEnabled      bool   `json:"namespaceConfusionProtectionEnabled"`
}

func (r *RepositoryDTO) ToApiRepository() *ApiRepository {
	var result ApiRepository
	copier.Copy(&result, r)
	return &result
}
