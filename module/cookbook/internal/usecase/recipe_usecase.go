package usecase

//go:generate mockgen -destination=../repository/mock/recipe_repo.go -source=recipe_usecase.go -package=mock RecipeRepository
//go:generate mockgen -destination=../repository/mock/recipe_ingredient_repo.go -source=recipe_usecase.go -package=mock RecipeIngredientRepository

import (
	"context"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
)

type ListRecipesFiter struct {
	CategoryID   uint64
	IngredientID uint64
}

type BulkRecipeIngredientParams []RecipeIngredientParams

type CreateRecipeParams struct {
	RecipeParams
	Ingredients BulkRecipeIngredientParams
}

type RecipeParams struct {
	Name        string
	Description string
	CategoryID  uint64
	Actor       string
}

type RecipeIngredientParams struct {
	Amount             float64
	IngredientID       uint64
	IngredientName     string
	IngredientUnitName string
	OrderingIndex      int
	Notes              string
	Actor              string
}

// RecipeRepository defines contract for recipe repository dependency
type RecipeRepository interface {
	Create(ctx context.Context, params CreateRecipeParams) (*entity.Recipe, error)
	Update(ctx context.Context, id uint64, params RecipeParams) (*entity.Recipe, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, filter ListRecipesFiter, limit, offset int) (entity.Recipes, error)
	GetSummary(ctx context.Context, id uint64) (entity.RecipeSummary, error)
}

// RecipeIngredientRepository defines contract for recipe ingredient repository dependency
type RecipeIngredientRepository interface {
	BulkCreate(ctx context.Context, recipeID uint64, params BulkRecipeIngredientParams) error
	Update(ctx context.Context, id uint64, params RecipeIngredientParams) (*entity.RecipeIngredient, error)
	Delete(ctx context.Context, id uint64) error
}

// RecipeUsecase is our recipe usecase object
type RecipeUsecase struct {
	recipeRepo           RecipeRepository
	recipeIngredientRepo RecipeIngredientRepository
}

// NewRecipeUsecase instantiates RecipeUsecase
func NewRecipeUsecase(recipeRepo RecipeRepository, recipeIngredientRepo RecipeIngredientRepository) *RecipeUsecase {
	return &RecipeUsecase{
		recipeRepo:           recipeRepo,
		recipeIngredientRepo: recipeIngredientRepo,
	}
}

// CreateRecipe creates a new recipe
func (u *RecipeUsecase) CreateRecipe(ctx context.Context, params CreateRecipeParams) error {
	recipe, err := u.recipeRepo.Create(ctx, params)
	if err != nil {
		return err
	}

	return u.recipeIngredientRepo.BulkCreate(ctx, recipe.ID, params.Ingredients)
}

func (u *RecipeUsecase) BulkCreateRecipeIngredients(ctx context.Context, recipeID uint64, params BulkRecipeIngredientParams) error {
	return u.recipeIngredientRepo.BulkCreate(ctx, recipeID, params)
}

// UpdateRecipe updates a voyage
func (u *RecipeUsecase) UpdateRecipe(ctx context.Context, id uint64, params RecipeParams) (*entity.Recipe, error) {
	return u.recipeRepo.Update(ctx, id, params)
}

// UpdateRecipeIngredient updates a voyage
func (u *RecipeUsecase) UpdateRecipeIngredient(ctx context.Context, id uint64, params RecipeIngredientParams) (*entity.RecipeIngredient, error) {
	return u.recipeIngredientRepo.Update(ctx, id, params)
}

// DeleteRecipe updates a voyage
func (u *RecipeUsecase) DeleteRecipe(ctx context.Context, id uint64) error {
	return u.recipeRepo.Delete(ctx, id)
}

// DeleteRecipeIngredient updates a voyage
func (u *RecipeUsecase) DeleteRecipeIngredient(ctx context.Context, id uint64) error {
	return u.recipeIngredientRepo.Delete(ctx, id)
}

// ListRecipes retrieves a list of tracking
func (u *RecipeUsecase) ListRecipes(ctx context.Context, filter ListRecipesFiter, limit, offset int) (entity.Recipes, error) {
	lim := defaultLimit
	ofs := defaultOffset

	if limit > 0 {
		lim = limit
	}

	if offset > 0 {
		ofs = offset
	}

	return u.recipeRepo.List(ctx, filter, lim, ofs)
}

func (u *RecipeUsecase) GetRecipeSummary(ctx context.Context, id uint64) (entity.RecipeSummary, error) {
	return u.recipeRepo.GetSummary(ctx, id)
}
