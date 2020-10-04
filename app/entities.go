package app

import "time"

// PrincipalType type of principal
type PrincipalType int

// principal types
const (
	PrincipalUSER PrincipalType = iota
	PrincipalTEST
)

// Principal - entity that is performing the op
type Principal struct {
	Type PrincipalType
	ID   string
}

// Entry - text entry
type Entry struct {
	ID        int
	Text      string
	CreatorID string
	CreatedAt time.Time
	UpdatedAt time.Time
}
