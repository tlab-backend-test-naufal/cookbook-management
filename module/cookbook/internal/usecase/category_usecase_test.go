package usecase_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/repository/mock"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
	"testing"
)

func TestNewCategoryUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)

	categoryRepo := mock.NewMockCategoryRepository(ctrl)
	uc := usecase.NewCategoryUsecase(categoryRepo)

	assert.NotEmpty(t, uc)
}
