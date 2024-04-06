package entity

import "gorm.io/gorm"

type Repository struct {
	gorm.Model
	InstanceID string `gorm:"uniqueIndex:instance_name"`
	RepositoryDTO
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
