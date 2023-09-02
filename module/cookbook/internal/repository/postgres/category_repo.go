package postgres_repo

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"

	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
)

// CategoryPostgresRepository is the PostgreSQL implementation for CategoryRepository interface
type CategoryPostgresRepository struct {
	db *sqlx.DB
}

// NewCategoryPostgresRepository instantiates CategoryPostgresRepository
func NewCategoryPostgresRepository(db *sqlx.DB) *CategoryPostgresRepository {
	return &CategoryPostgresRepository{db: db}
}

type categoryDto struct {
	ID        uint64      `db:"id"`
	Name      string      `db:"name"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy string      `db:"created_by"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy null.String `db:"updated_by"`
	IsDeleted bool        `db:"is_deleted"`
}

func (c categoryDto) toEntity() *entity.Category {
	return &entity.Category{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		CreatedBy: c.CreatedBy,
		UpdatedAt: c.UpdatedAt,
		UpdatedBy: c.UpdatedBy,
		IsDeleted: c.IsDeleted,
	}
}

const selectCategoryQuery = `
select id, name, created_at, created_by, updated_at, updated_by from categories
where is_deleted = false
limit $1 offset $2;
`

// List retrieves a list of categories with offset and limit
func (r *CategoryPostgresRepository) List(ctx context.Context, limit, offset int) (res entity.Categories, err error) {
	var dtos []categoryDto

	err = r.db.SelectContext(ctx, &dtos, selectCategoryQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, dto := range dtos {
		res = append(res, dto.toEntity())
	}

	return res, nil
}

const insertCategoryQuery = `
INSERT INTO categories (name, created_at, created_by)
VALUES ($1, $2, $3) RETURNING id
`

// Create creates a new category
func (r *CategoryPostgresRepository) Create(ctx context.Context, params usecase.CategoryParams) (*entity.Category, error) {
	dto := categoryDtoForCreate(params)

	err := r.db.QueryRowxContext(ctx, insertCategoryQuery, dto.Name, dto.CreatedAt, dto.CreatedBy).Scan(&dto.ID)
	if err != nil {
		return nil, err
	}

	return &entity.Category{
		ID:        dto.ID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		CreatedBy: dto.CreatedBy,
		UpdatedAt: dto.UpdatedAt,
		UpdatedBy: dto.UpdatedBy,
	}, nil
}

// Update updates a category by its ID
func (r *CategoryPostgresRepository) Update(ctx context.Context, id uint64, params usecase.CategoryParams) (*entity.Category, error) {
	dto, query := categoryDtoForUpdate(id, params, nil)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return nil, entity.ErrCategoryNotFound
	}

	if err != nil {
		return nil, err
	}

	return &entity.Category{
		ID:        dto.ID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}, nil
}

// Delete deletes a category by its ID
func (r *CategoryPostgresRepository) Delete(ctx context.Context, id uint64) error {
	dto, query := categoryDtoForDelete(id)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return entity.ErrCategoryNotFound
	}

	return err
}

func categoryDtoForCreate(params usecase.CategoryParams) categoryDto {
	return categoryDto{
		Name:      params.Name,
		CreatedAt: time.Now(),
		CreatedBy: params.Actor,
	}
}

func categoryDtoForUpdate(id uint64, params usecase.CategoryParams, isDeleted *bool) (dto categoryDto, query string) {
	var qb strings.Builder

	qb.WriteString("UPDATE categories SET ")

	if params.Name != "" {
		qb.WriteString("name = :name, ")
		dto.Name = params.Name
	}

	if isDeleted != nil {
		qb.WriteString("is_deleted = :is_deleted, ")
		dto.IsDeleted = *isDeleted
	}

	qb.WriteString("updated_at = :updated_at, ")
	dto.UpdatedAt = null.TimeFrom(time.Now())

	qb.WriteString("updated_by = :updated_by ")
	dto.UpdatedBy = null.StringFrom(params.Actor)

	qb.WriteString("WHERE id = :id")
	dto.ID = id

	return dto, qb.String()
}

func categoryDtoForDelete(id uint64) (dto categoryDto, query string) {
	isDeleted := true
	return categoryDtoForUpdate(id, usecase.CategoryParams{}, &isDeleted)
}
