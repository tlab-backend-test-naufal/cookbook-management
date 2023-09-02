package config

import (
	"github.com/jmoiron/sqlx"
	cookbookPostgresRepo "github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/repository/postgres"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
	cookbookRest "github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/rest"
)

func RegisterCookbookHandler(db *sqlx.DB) *cookbookRest.CookbookHandler {
	categoryRepo := cookbookPostgresRepo.NewCategoryPostgresRepository(db)

	ingredientRepo := cookbookPostgresRepo.NewIngredientPostgresRepository(db)
	ingredientUnitRepo := cookbookPostgresRepo.NewIngredientUnitPostgresRepository(db)

	recipeRepo := cookbookPostgresRepo.NewRecipePostgresRepository(db)
	recipeIngredientRepo := cookbookPostgresRepo.NewRecipeIngredientPostgresRepository(db)

	cookbookUc := usecase.NewCategoryUsecase(categoryRepo)
	ingredientUc := usecase.NewIngredientUsecase(ingredientRepo, ingredientUnitRepo)

	recipeUc := usecase.NewRecipeUsecase(recipeRepo, recipeIngredientRepo)

	return cookbookRest.NewCookbookHandler(cookbookUc, ingredientUc, recipeUc)
}
