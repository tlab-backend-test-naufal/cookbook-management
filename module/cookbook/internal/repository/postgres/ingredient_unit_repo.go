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

// IngredientUnitPostgresRepository is the PostgreSQL implementation for IngredientUnitRepository interface
type IngredientUnitPostgresRepository struct {
	db *sqlx.DB
}

// NewIngredientUnitPostgresRepository instantiates IngredientUnitPostgresRepository
func NewIngredientUnitPostgresRepository(db *sqlx.DB) *IngredientUnitPostgresRepository {
	return &IngredientUnitPostgresRepository{db: db}
}

type ingredientUnitDto struct {
	ID        uint64      `db:"id"`
	Name      string      `db:"name"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy string      `db:"created_by"`
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy null.String `db:"updated_by"`
	IsDeleted bool        `db:"is_deleted"`
}

func (c ingredientUnitDto) toEntity() *entity.IngredientUnit {
	return &entity.IngredientUnit{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		CreatedBy: c.CreatedBy,
		UpdatedAt: c.UpdatedAt,
		UpdatedBy: c.UpdatedBy,
		IsDeleted: c.IsDeleted,
	}
}

const selectIngredientUnitQuery = `
select id, name, created_at, created_by, updated_at, updated_by, is_deleted from ingredient_units
where is_deleted = false
limit $1 offset $2;
`

// List retrieves a list of ingredient units with offset and limit
func (r *IngredientUnitPostgresRepository) List(ctx context.Context, limit, offset int) (res entity.IngredientUnits, err error) {
	var dtos []ingredientUnitDto

	err = r.db.SelectContext(ctx, &dtos, selectIngredientUnitQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, dto := range dtos {
		res = append(res, dto.toEntity())
	}

	return res, nil
}

const insertIngredientUnitQuery = `
INSERT INTO ingredient_units (name, created_at, created_by)
VALUES ($1, $2, $3) RETURNING id
`

// Create creates a new ingredientUnit
func (r *IngredientUnitPostgresRepository) Create(ctx context.Context, params usecase.IngredientUnitParams) (*entity.IngredientUnit, error) {
	dto := ingredientUnitDtoForCreate(params)

	err := r.db.QueryRowxContext(ctx, insertIngredientUnitQuery, dto.Name, dto.CreatedAt, dto.CreatedBy).Scan(&dto.ID)
	if err != nil {
		return nil, err
	}

	return &entity.IngredientUnit{
		ID:        dto.ID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		CreatedBy: dto.CreatedBy,
		UpdatedAt: dto.UpdatedAt,
		UpdatedBy: dto.UpdatedBy,
	}, nil
}

// Update updates a ingredientUnit by its ID
func (r *IngredientUnitPostgresRepository) Update(ctx context.Context, id uint64, params usecase.IngredientUnitParams) (*entity.IngredientUnit, error) {
	dto, query := ingredientUnitDtoForUpdate(id, params, nil)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return nil, entity.ErrIngredientUnitNotFound
	}

	if err != nil {
		return nil, err
	}

	return &entity.IngredientUnit{
		ID:        dto.ID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		CreatedBy: dto.CreatedBy,
		UpdatedBy: dto.UpdatedBy,
	}, nil
}

// Delete deletes a ingredientUnit by its ID
func (r *IngredientUnitPostgresRepository) Delete(ctx context.Context, id uint64) error {
	dto, query := ingredientUnitDtoForDelete(id)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return entity.ErrIngredientUnitNotFound
	}

	return err
}

func ingredientUnitDtoForCreate(params usecase.IngredientUnitParams) ingredientUnitDto {
	return ingredientUnitDto{
		Name:      params.Name,
		CreatedAt: time.Now(),
		CreatedBy: params.Actor,
	}
}

func ingredientUnitDtoForUpdate(id uint64, params usecase.IngredientUnitParams, isDeleted *bool) (dto ingredientUnitDto, query string) {
	var qb strings.Builder

	qb.WriteString("UPDATE ingredient_units SET ")

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

func ingredientUnitDtoForDelete(id uint64) (dto ingredientUnitDto, query string) {
	isDeleted := true
	return ingredientUnitDtoForUpdate(id, usecase.IngredientUnitParams{}, &isDeleted)
}
