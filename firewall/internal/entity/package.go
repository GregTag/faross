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
	FinalScore float64
	Report     string
}

func (s State) ToSring() string {
	switch s {
	case Healthy:
		return "healthy"
	case Quarantined:
		return "quarantined"
	case Unquarantined:
		return "unquarantined"
	default:
		return "undefined"
	}
}
