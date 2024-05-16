package entity

import "errors"

var (
	ErrPackageNotFound      = errors.New("package not found")
	ErrNothingUnquarantined = errors.New("no packages to unquarantine")
	ErrPending              = errors.New("package is pending")
)
