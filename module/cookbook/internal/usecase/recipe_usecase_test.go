package usecase_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/repository/mock"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
)

func TestNewRecipeUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)

	recipeRepo := mock.NewMockRecipeRepository(ctrl)
	recipeIngredientRepo := mock.NewMockRecipeIngredientRepository(ctrl)

	uc := usecase.NewRecipeUsecase(recipeRepo, recipeIngredientRepo)

	assert.NotEmpty(t, uc)
}
