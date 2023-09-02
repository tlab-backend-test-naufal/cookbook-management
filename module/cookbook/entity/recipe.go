package entity

import (
	"time"

	"github.com/guregu/null"
)

// Recipes is the plural form of Recipe
type Recipes []*Recipe

// Recipe holds our recipe entity
type Recipe struct {
	ID          uint64
	Name        string
	Description string
	CategoryID  uint64
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   null.Time
	UpdatedBy   null.String
	IsDeleted   bool
}

// RecipeSummary is a summary of a recipe with its ingredients
type RecipeSummary struct {
	Recipe
	Ingredients RecipeIngredients
}
