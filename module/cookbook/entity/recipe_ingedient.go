package entity

import (
	"time"

	"github.com/guregu/null"
)

// RecipeIngredients is the plural form of RecipeIngredient
type RecipeIngredients []*RecipeIngredient

// RecipeIngredient holds our recipe entity
type RecipeIngredient struct {
	ID                 uint64
	RecipeID           uint64
	IngredientID       uint64
	IngredientName     string
	IngredientUnitName string
	Amount             float64
	OrderingIndex      int
	Notes              string
	CreatedAt          time.Time
	CreatedBy          string
	UpdatedAt          null.Time
	UpdatedBy          null.String
	IsDeleted          bool
}
