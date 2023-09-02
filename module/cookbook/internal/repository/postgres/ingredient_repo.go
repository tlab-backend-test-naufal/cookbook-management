package postgres_repo

import (
	"context"
	"database/sql"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
	"strings"
	"time"
)

// IngredientPostgresRepository is the PostgreSQL implementation for IngredientRepository interface
type IngredientPostgresRepository struct {
	db *sqlx.DB
}

// NewIngredientPostgresRepository instantiates IngredientPostgresRepository
func NewIngredientPostgresRepository(db *sqlx.DB) *IngredientPostgresRepository {
	return &IngredientPostgresRepository{db: db}
}

type ingredientDto struct {
	ID        uint64      `db:"id"`
	Name      string      `db:"name"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy string      `db:"created_by"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy null.String `db:"updated_by"`
	IsDeleted bool        `db:"is_deleted"`
}

func (c ingredientDto) toEntity() *entity.Ingredient {
	return &entity.Ingredient{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		CreatedBy: c.CreatedBy,
		UpdatedAt: c.UpdatedAt,
		UpdatedBy: c.UpdatedBy,
		IsDeleted: c.IsDeleted,
	}
}

const selectIngredientQuery = `
select id, name, created_at, created_by, updated_at, updated_by, is_deleted from ingredients
where is_deleted = false
limit $1 offset $2;
`

// List retrieves a list of ingredients with offset and limit
func (r *IngredientPostgresRepository) List(ctx context.Context, limit, offset int) (res entity.Ingredients, err error) {
	var dtos []ingredientDto

	err = r.db.SelectContext(ctx, &dtos, selectIngredientQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, dto := range dtos {
		res = append(res, dto.toEntity())
	}

	return res, nil
}

const insertIngredientQuery = `
INSERT INTO ingredients (name, created_at, created_by)
VALUES ($1, $2, $3) RETURNING id
`

// Create creates a new ingredient
func (r *IngredientPostgresRepository) Create(ctx context.Context, params usecase.IngredientParams) (*entity.Ingredient, error) {
	dto := ingredientDtoForCreate(params)

	err := r.db.QueryRowxContext(ctx, insertIngredientQuery, dto.Name, dto.CreatedAt, dto.CreatedBy).Scan(&dto.ID)
	if err != nil {
		return nil, err
	}

	return &entity.Ingredient{
		ID:        dto.ID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		CreatedBy: dto.CreatedBy,
		UpdatedAt: dto.UpdatedAt,
		UpdatedBy: dto.UpdatedBy,
	}, nil
}

// Update updates a ingredient by its ID
func (r *IngredientPostgresRepository) Update(ctx context.Context, id uint64, params usecase.IngredientParams) (*entity.Ingredient, error) {
	dto, query := ingredientDtoForUpdate(id, params, nil)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return nil, entity.ErrIngredientNotFound
	}

	if err != nil {
		return nil, err
	}

	return &entity.Ingredient{
		ID:        dto.ID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}, nil
}

// Delete deletes a ingredient by its ID
func (r *IngredientPostgresRepository) Delete(ctx context.Context, id uint64) error {
	dto, query := ingredientDtoForDelete(id)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return entity.ErrIngredientNotFound
	}

	return err
}

func ingredientDtoForCreate(params usecase.IngredientParams) ingredientDto {
	return ingredientDto{
		Name:      params.Name,
		CreatedAt: time.Now(),
		CreatedBy: params.Actor,
	}
}

func ingredientDtoForUpdate(id uint64, params usecase.IngredientParams, isDeleted *bool) (dto ingredientDto, query string) {
	var qb strings.Builder

	qb.WriteString("UPDATE ingredients SET ")

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

func ingredientDtoForDelete(id uint64) (dto ingredientDto, query string) {
	isDeleted := true
	return ingredientDtoForUpdate(id, usecase.IngredientParams{}, &isDeleted)
}
