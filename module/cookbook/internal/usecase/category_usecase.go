package usecase

//go:generate mockgen -destination=../repository/mock/category_repository.go -source=category_usecase.go -package=mock CategoryRepository

import (
	"context"

	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
)

type CategoryParams struct {
	Name  string
	Actor string
}

// CategoryRepository defines contract for ingredient repository dependency
type CategoryRepository interface {
	Create(ctx context.Context, params CategoryParams) (*entity.Category, error)
	Update(ctx context.Context, id uint64, params CategoryParams) (*entity.Category, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, limit, offset int) (entity.Categories, error)
}

// CategoryUsecase is our ingredient usecase object
type CategoryUsecase struct {
	categoryRepo CategoryRepository
}

// NewCategoryUsecase instantiates CategoryUsecase
func NewCategoryUsecase(categoryRepo CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory creates a new category
func (u *CategoryUsecase) CreateCategory(ctx context.Context, params CategoryParams) (*entity.Category, error) {
	return u.categoryRepo.Create(ctx, params)
}

// UpdateCategory updates a category
func (u *CategoryUsecase) UpdateCategory(ctx context.Context, id uint64, params CategoryParams) (*entity.Category, error) {
	return u.categoryRepo.Update(ctx, id, params)
}

// DeleteCategory deletes a category
func (u *CategoryUsecase) DeleteCategory(ctx context.Context, id uint64) error {
	return u.categoryRepo.Delete(ctx, id)
}

// ListCategories retrieves a list of categories
func (u *CategoryUsecase) ListCategories(ctx context.Context, limit, offset int) (entity.Categories, error) {
	lim := defaultLimit
	ofs := defaultOffset

	if limit > 0 {
		lim = limit
	}

	if offset > 0 {
		ofs = offset
	}

	return u.categoryRepo.List(ctx, lim, ofs)
}
