package postgres_repo

import (
	"context"
	"database/sql"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"

	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
)

// RecipePostgresRepository is the PostgreSQL implementation for RecipeRepository interface
type RecipePostgresRepository struct {
	db *sqlx.DB
}

// NewRecipePostgresRepository instantiates RecipePostgresRepository
func NewRecipePostgresRepository(db *sqlx.DB) *RecipePostgresRepository {
	return &RecipePostgresRepository{db: db}
}

type recipeDto struct {
	ID          uint64      `db:"id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
	CategoryID  uint64      `db:"category_id"`
	CreatedAt   time.Time   `db:"created_at"`
	CreatedBy   string      `db:"created_by"`
	UpdatedAt   null.Time   `db:"updated_at"`
	UpdatedBy   null.String `db:"updated_by"`
	IsDeleted   bool        `db:"is_deleted"`
}

type recipeSummaryDto struct {
	ID                 uint64      `db:"id"`
	Name               string      `db:"name"`
	Description        string      `db:"description"`
	CategoryID         uint64      `db:"category_id"`
	RecipeIngredientID uint64      `db:"recipe_ingredient_id"`
	IngredientName     string      `db:"ingredient_name"`
	IngredientUnitName string      `db:"ingredient_unit_name"`
	Amount             float64     `db:"amount"`
	Notes              string      `db:"notes"`
	OrderingIndex      int         `db:"ordering_index"`
	CreatedAt          time.Time   `db:"created_at"`
	CreatedBy          string      `db:"created_by"`
	UpdatedAt          null.Time   `db:"updated_at"`
	UpdatedBy          null.String `db:"updated_by"`
	IsDeleted          bool        `db:"is_deleted"`
}

func (c recipeDto) toEntity() *entity.Recipe {
	return &entity.Recipe{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CategoryID:  c.CategoryID,
		CreatedAt:   c.CreatedAt,
		CreatedBy:   c.CreatedBy,
		UpdatedAt:   c.UpdatedAt,
		UpdatedBy:   c.UpdatedBy,
		IsDeleted:   c.IsDeleted,
	}
}

const selectRecipeQuery = `
select
       r.id as id,
       r.name,
       r.description,
       r.category_id,
       r.created_at,
       r.created_by,
       r.updated_at,
       r.updated_by
from recipes r
left join recipe_ingredients ri on r.id = ri.recipe_id
where r.is_deleted = false`

// List retrieves a list of categories with offset and limit
func (r *RecipePostgresRepository) List(ctx context.Context, filter usecase.ListRecipesFiter, limit, offset int) (res entity.Recipes, err error) {
	var dtos []recipeDto

	query := selectRecipeQuery

	args := []interface{}{limit, offset}

	if filter.CategoryID > 0 {
		query += "\nand r.category_id = $3"
		args = append(args, filter.CategoryID)
	}

	if filter.IngredientID > 0 {
		query += "\nand ri.ingredient_id = $4"
		args = append(args, filter.IngredientID)
	}

	query += "\nlimit $1 offset $2;"

	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, err
	}

	for _, dto := range dtos {
		res = append(res, dto.toEntity())
	}

	return res, nil
}

const insertRecipeQuery = `
INSERT INTO recipes (name, description, category_id, created_at, created_by)
VALUES ($1, $2, $3, $4, $5) RETURNING id
`

// Create creates a new Recipe
func (r *RecipePostgresRepository) Create(ctx context.Context, params usecase.CreateRecipeParams) (*entity.Recipe, error) {
	dto := recipeDtoForCreate(params)

	err := r.db.QueryRowxContext(ctx, insertRecipeQuery, dto.Name, dto.Description, dto.CategoryID, dto.CreatedAt, dto.CreatedBy).Scan(&dto.ID)
	if err != nil {
		return nil, err
	}

	return &entity.Recipe{
		ID:        dto.ID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		CreatedBy: dto.CreatedBy,
		UpdatedAt: dto.UpdatedAt,
		UpdatedBy: dto.UpdatedBy,
	}, nil
}

// Update updates a Recipe by its ID
func (r *RecipePostgresRepository) Update(ctx context.Context, id uint64, params usecase.RecipeParams) (*entity.Recipe, error) {
	dto, query := recipeDtoForUpdate(id, params, nil)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return nil, entity.ErrRecipeNotFound
	}

	if err != nil {
		return nil, err
	}

	return &entity.Recipe{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		CategoryID:  dto.CategoryID,
		CreatedAt:   dto.CreatedAt,
		CreatedBy:   dto.CreatedBy,
		UpdatedAt:   dto.UpdatedAt,
		UpdatedBy:   dto.UpdatedBy,
		IsDeleted:   dto.IsDeleted,
	}, nil
}

// Delete deletes a Recipe by its ID
func (r *RecipePostgresRepository) Delete(ctx context.Context, id uint64) error {
	dto, query := recipeDtoForDelete(id)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return entity.ErrRecipeNotFound
	}

	return err
}

const selectRecipeSummaryQuery = `
select
       r.id as id,
       r.name,
       r.description,
       r.category_id,
       ri.id as recipe_ingredient_id,
       ri.ingredient_name as ingredient_name,
       ri.ingredient_unit_name as ingredient_unit_name,
       ri.amount as amount,
       ri.notes as notes,
       ri.ordering_index as ordering_index,
       r.created_at,
       r.created_by,
       r.updated_at,
       r.updated_by
from recipes r
left join recipe_ingredients ri on r.id = ri.recipe_id
where r.is_deleted = false
and r.id = $1
order by ordering_index;
`

func (r *RecipePostgresRepository) GetSummary(ctx context.Context, id uint64) (entity.RecipeSummary, error) {
	var dtos []recipeSummaryDto

	err := r.db.SelectContext(ctx, &dtos, selectRecipeSummaryQuery, id)
	if err != nil {
		return entity.RecipeSummary{}, err
	}

	if len(dtos) == 0 {
		return entity.RecipeSummary{}, nil
	}

	return constructRecipeSummary(dtos), nil
}

func constructRecipeSummary(dtos []recipeSummaryDto) entity.RecipeSummary {
	var ingredients entity.RecipeIngredients

	for _, dto := range dtos {
		ingredients = append(ingredients, &entity.RecipeIngredient{
			ID:                 dto.ID,
			IngredientName:     dto.IngredientName,
			IngredientUnitName: dto.IngredientUnitName,
			Amount:             dto.Amount,
			OrderingIndex:      dto.OrderingIndex,
			Notes:              dto.Notes,
			CreatedAt:          dto.CreatedAt,
			CreatedBy:          dto.CreatedBy,
			UpdatedAt:          dto.UpdatedAt,
			UpdatedBy:          dto.UpdatedBy,
			IsDeleted:          dto.IsDeleted,
		})
	}

	return entity.RecipeSummary{
		Recipe: entity.Recipe{
			ID:          dtos[0].ID,
			Name:        dtos[0].Name,
			Description: dtos[0].Description,
			CategoryID:  dtos[0].CategoryID,
			CreatedAt:   dtos[0].CreatedAt,
			CreatedBy:   dtos[0].CreatedBy,
			UpdatedAt:   dtos[0].UpdatedAt,
			UpdatedBy:   dtos[0].UpdatedBy,
			IsDeleted:   dtos[0].IsDeleted,
		},
		Ingredients: ingredients,
	}
}

func recipeDtoForCreate(params usecase.CreateRecipeParams) recipeDto {
	return recipeDto{
		Name:        params.Name,
		Description: params.Description,
		CategoryID:  params.CategoryID,
		CreatedAt:   time.Now(),
		CreatedBy:   params.Actor,
	}
}

func recipeDtoForUpdate(id uint64, params usecase.RecipeParams, isDeleted *bool) (dto recipeDto, query string) {
	var qb strings.Builder

	qb.WriteString("UPDATE recipes SET ")

	if params.Name != "" {
		qb.WriteString("name = :name, ")
		dto.Name = params.Name
	}

	if params.Description != "" {
		qb.WriteString("description = :description, ")
		dto.Name = params.Description
	}

	if params.CategoryID != 0 {
		qb.WriteString("category_id = :category_id, ")
		dto.Name = params.Description
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

func recipeDtoForDelete(id uint64) (dto recipeDto, query string) {
	isDeleted := true
	return recipeDtoForUpdate(id, usecase.RecipeParams{}, &isDeleted)
}
