package rest

import (
	"context"

	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
)

// CategoryUsecase defines the contract for cookbook usecase dependency
type CategoryUsecase interface {
	CreateCategory(ctx context.Context, params usecase.CategoryParams) (*entity.Category, error)
	UpdateCategory(ctx context.Context, id uint64, params usecase.CategoryParams) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
	ListCategories(ctx context.Context, limit, offset int) (entity.Categories, error)
}

// IngredientUsecase defines the contract for voyage usecase dependency
type IngredientUsecase interface {
	CreateIngredient(ctx context.Context, params usecase.IngredientParams) (*entity.Ingredient, error)
	UpdateIngredient(ctx context.Context, id uint64, params usecase.IngredientParams) (*entity.Ingredient, error)
	DeleteIngredient(ctx context.Context, id uint64) error
	ListIngredients(ctx context.Context, limit, offset int) (entity.Ingredients, error)
	CreateIngredientUnit(ctx context.Context, params usecase.IngredientUnitParams) (*entity.IngredientUnit, error)
	UpdateIngredientUnit(ctx context.Context, id uint64, params usecase.IngredientUnitParams) (*entity.IngredientUnit, error)
	DeleteIngredientUnit(ctx context.Context, id uint64) error
	ListIngredientUnits(ctx context.Context, limit, offset int) (entity.IngredientUnits, error)
}

type RecipeUsecase interface {
	CreateRecipe(ctx context.Context, params usecase.CreateRecipeParams) error
	UpdateRecipe(ctx context.Context, id uint64, params usecase.RecipeParams) (*entity.Recipe, error)
	DeleteRecipe(ctx context.Context, id uint64) error
	BulkCreateRecipeIngredients(ctx context.Context, recipeID uint64, params usecase.BulkRecipeIngredientParams) error
	UpdateRecipeIngredient(ctx context.Context, id uint64, params usecase.RecipeIngredientParams) (*entity.RecipeIngredient, error)
	DeleteRecipeIngredient(ctx context.Context, id uint64) error
	ListRecipes(ctx context.Context, filter usecase.ListRecipesFiter, limit, offset int) (entity.Recipes, error)
	GetRecipeSummary(ctx context.Context, id uint64) (entity.RecipeSummary, error)
}

// CookbookHandler is our GraphQL resolver object
type CookbookHandler struct {
	categoryUsecase   CategoryUsecase
	ingredientUsecase IngredientUsecase
	recipeUsecase     RecipeUsecase
}

// NewCookbookHandler instantiates cookbookHandler
func NewCookbookHandler(categoryUsecase CategoryUsecase, ingredientUsecase IngredientUsecase, recipeUsecase RecipeUsecase) *CookbookHandler {
	return &CookbookHandler{
		categoryUsecase:   categoryUsecase,
		ingredientUsecase: ingredientUsecase,
		recipeUsecase:     recipeUsecase,
	}
}
