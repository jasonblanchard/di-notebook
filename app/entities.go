package app

import "time"

// PrincipalType type of principle
type PrincipalType int

// principle types
const (
	PrincipleUSER PrincipalType = iota
	PrincipleTEST
)

// Principle - entity that is performing the op
type Principle struct {
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
