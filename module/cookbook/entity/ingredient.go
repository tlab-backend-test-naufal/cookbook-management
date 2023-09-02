package entity

import (
	"time"

	"github.com/guregu/null"
)

// Ingredients is the plural form of Ingredient
type Ingredients []*Ingredient

// Ingredient holds our ingredient entity
type Ingredient struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt null.Time
	UpdatedBy null.String
	IsDeleted bool
}
