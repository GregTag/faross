package entity

import (
	"gorm.io/gorm"
)

type State int

const (
	Undefined State = iota
	Pending
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
	Comment    string
}

func (s State) ToSring() string {
	switch s {
	case Pending:
		return "pending"
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
