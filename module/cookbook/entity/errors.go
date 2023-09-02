package entity

import "github.com/tlab-backend-test-naufal/cookbook-management/internal/liberr"

var (
	ErrCategoryNotFound       = liberr.NewErrorDetails("COOKBOOK_COOKBOOK-MANAGEMENT_CATEGORY-NOT-FOUND", "Category is not found")
	ErrIngredientNotFound     = liberr.NewErrorDetails("COOKBOOK_COOKBOOK-MANAGEMENT_INGREDIENT-NOT-FOUND", "Ingredient is not found")
	ErrIngredientUnitNotFound = liberr.NewErrorDetails("COOKBOOK_COOKBOOK-MANAGEMENT_INGREDIENT-UNIT-NOT-FOUND", "Ingredient unit is not found")
	ErrRecipeNotFound         = liberr.NewErrorDetails("COOKBOOK_COOKBOOK-MANAGEMENT_RECIPE-UNIT-NOT-FOUND", "Recipe is not found")
)
