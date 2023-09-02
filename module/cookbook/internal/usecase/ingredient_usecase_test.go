package usecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/repository/mock"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
)

func TestNewIngredientUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)

	ingredientRepo := mock.NewMockIngredientRepository(ctrl)
	ingredientUnitRepo := mock.NewMockIngredientUnitRepository(ctrl)

	uc := usecase.NewIngredientUsecase(ingredientRepo, ingredientUnitRepo)

	assert.NotEmpty(t, uc)
}
