package entity

import (
	"github.com/guregu/null"
	"time"
)

// IngredientUnits is the plural form of IngredientUnit
type IngredientUnits []*IngredientUnit

// IngredientUnit holds our ingredient unit entity
type IngredientUnit struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt null.Time
	UpdatedBy null.String
	IsDeleted bool
}
