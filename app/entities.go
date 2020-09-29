package app

import "time"

// PrincipleType type of principle
type PrincipleType int

// principle types
const (
	PrincipleTypeUser PrincipleType = iota
	PrincipleTypeTest
)

// Principle - entity that is performing the op
type Principle struct {
	Type PrincipleType
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
