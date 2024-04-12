package entity

import (
	"gorm.io/gorm"
)

type State int

const (
	Undefined State = iota
	Healthy
	Quarantined
	Unquarantined
)

type Package struct {
	gorm.Model
	Pathname   string
	Purl       string `gorm:"unique"`
	State      State
	FinalScore float32
	Report     string
}
