package usecase

//go:generate mockgen -destination=../repository/mock/ingredient_repo.go -source=ingredient_usecase.go -package=mock IngredientRepository
//go:generate mockgen -destination=../repository/mock/ingredient_unit_repo.go -source=ingredient_usecase.go -package=mock IngredientUnitRepository

import (
	"context"

	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
)

type IngredientParams struct {
	Name  string
	Actor string
}

type IngredientUnitParams struct {
	Name  string
	Actor string
}

// IngredientRepository defines contract for ingredient repository dependency
type IngredientRepository interface {
	Create(ctx context.Context, params IngredientParams) (*entity.Ingredient, error)
	Update(ctx context.Context, id uint64, params IngredientParams) (*entity.Ingredient, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, limit, offset int) (entity.Ingredients, error)
}

// IngredientUnitRepository defines contract for ingredient unit repository dependency
type IngredientUnitRepository interface {
	Create(ctx context.Context, params IngredientUnitParams) (*entity.IngredientUnit, error)
	Update(ctx context.Context, id uint64, params IngredientUnitParams) (*entity.IngredientUnit, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, limit, offset int) (entity.IngredientUnits, error)
}

// IngredientUsecase is our ingredient usecase object
type IngredientUsecase struct {
	ingredientRepo     IngredientRepository
	ingredientUnitRepo IngredientUnitRepository
}

// NewIngredientUsecase instantiates IngredientUsecase
func NewIngredientUsecase(ingredientRepo IngredientRepository, ingredientUnitRepo IngredientUnitRepository) *IngredientUsecase {
	return &IngredientUsecase{
		ingredientRepo:     ingredientRepo,
		ingredientUnitRepo: ingredientUnitRepo,
	}
}

// CreateIngredient creates a new Ingredient
func (u *IngredientUsecase) CreateIngredient(ctx context.Context, params IngredientParams) (*entity.Ingredient, error) {
	return u.ingredientRepo.Create(ctx, params)
}

// UpdateIngredient updates a Ingredient
func (u *IngredientUsecase) UpdateIngredient(ctx context.Context, id uint64, params IngredientParams) (*entity.Ingredient, error) {
	return u.ingredientRepo.Update(ctx, id, params)
}

// DeleteIngredient deletes a Ingredient
func (u *IngredientUsecase) DeleteIngredient(ctx context.Context, id uint64) error {
	return u.ingredientRepo.Delete(ctx, id)
}

// ListIngredients retrieves a list of Ingredients
func (u *IngredientUsecase) ListIngredients(ctx context.Context, limit, offset int) (entity.Ingredients, error) {
	lim := defaultLimit
	ofs := defaultOffset

	if limit > 0 {
		lim = limit
	}

	if offset > 0 {
		ofs = offset
	}

	return u.ingredientRepo.List(ctx, lim, ofs)
}

// CreateIngredientUnit creates a new Ingredient
func (u *IngredientUsecase) CreateIngredientUnit(ctx context.Context, params IngredientUnitParams) (*entity.IngredientUnit, error) {
	return u.ingredientUnitRepo.Create(ctx, params)
}

// UpdateIngredientUnit updates a Ingredient
func (u *IngredientUsecase) UpdateIngredientUnit(ctx context.Context, id uint64, params IngredientUnitParams) (*entity.IngredientUnit, error) {
	return u.ingredientUnitRepo.Update(ctx, id, params)
}

// DeleteIngredientUnit deletes a Ingredient
func (u *IngredientUsecase) DeleteIngredientUnit(ctx context.Context, id uint64) error {
	return u.ingredientUnitRepo.Delete(ctx, id)
}

// ListIngredientUnits retrieves a list of Ingredients
func (u *IngredientUsecase) ListIngredientUnits(ctx context.Context, limit, offset int) (entity.IngredientUnits, error) {
	lim := defaultLimit
	ofs := defaultOffset

	if limit > 0 {
		lim = limit
	}

	if offset > 0 {
		ofs = offset
	}

	return u.ingredientUnitRepo.List(ctx, lim, ofs)
}
