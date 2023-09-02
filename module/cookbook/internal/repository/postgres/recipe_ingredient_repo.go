package postgres_repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
	"strings"
	"time"
)

type recipeIngredientDto struct {
	ID                 uint64      `db:"id"`
	RecipeID           uint64      `db:"recipe_id"`
	IngredientID       uint64      `db:"ingredient_id"`
	IngredientName     string      `db:"ingredient_name"`
	IngredientUnitName string      `db:"ingredient_unit_name"`
	Amount             float64     `db:"amount"`
	OrderingIndex      int         `db:"ordering_index"`
	Notes              string      `db:"notes"`
	CreatedAt          time.Time   `db:"created_at"`
	CreatedBy          string      `db:"created_by"`
	UpdatedAt          null.Time   `db:"updated_at"`
	UpdatedBy          null.String `db:"updated_by"`
	IsDeleted          bool        `db:"is_deleted"`
}

// RecipeIngredientPostgresRepository is the PostgreSQL implementation for RecipeIngredientRepository interface
type RecipeIngredientPostgresRepository struct {
	db *sqlx.DB
}

// NewRecipeIngredientPostgresRepository instantiates RecipeIngredientPostgresRepository
func NewRecipeIngredientPostgresRepository(db *sqlx.DB) *RecipeIngredientPostgresRepository {
	return &RecipeIngredientPostgresRepository{db: db}
}

const bulkInsertRecipeIngredientsQuery = `
INSERT INTO recipe_ingredients (recipe_id, ingredient_id, ingredient_name, ingredient_unit_name, amount, ordering_index, notes, created_at, created_by)
`

func (r *RecipeIngredientPostgresRepository) BulkCreate(ctx context.Context, recipeID uint64, params usecase.BulkRecipeIngredientParams) error {
	var args []interface{}
	var count int

	var dollars []string
	var values []string

	var rowLen int

	for _, p := range params {
		var argsTmp []interface{}

		argsTmp = append(argsTmp, recipeID)
		argsTmp = append(argsTmp, p.IngredientID)
		argsTmp = append(argsTmp, p.IngredientName)
		argsTmp = append(argsTmp, p.IngredientUnitName)
		argsTmp = append(argsTmp, p.Amount)
		argsTmp = append(argsTmp, p.OrderingIndex)
		argsTmp = append(argsTmp, p.Notes)
		argsTmp = append(argsTmp, time.Now())
		argsTmp = append(argsTmp, p.Actor)

		args = append(args, argsTmp...)
		rowLen = len(argsTmp)

		count += rowLen
	}

	for i := 1; i <= len(args); i++ {
		dollars = append(dollars, fmt.Sprintf("$%d", i))
		if i%rowLen == 0 {
			values = append(values, fmt.Sprintf("(%s)", strings.Join(dollars, ",")))
			dollars = []string{}
		}
	}

	query := fmt.Sprintf("%s VALUES %s", bulkInsertRecipeIngredientsQuery, strings.Join(values, ","))

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *RecipeIngredientPostgresRepository) Update(ctx context.Context, id uint64, params usecase.RecipeIngredientParams) (*entity.RecipeIngredient, error) {
	dto, query := recipeIngredientDtoForUpdate(id, params, nil)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return nil, entity.ErrRecipeNotFound
	}

	if err != nil {
		return nil, err
	}

	return &entity.RecipeIngredient{
		ID:                 dto.ID,
		RecipeID:           dto.RecipeID,
		IngredientID:       dto.IngredientID,
		IngredientName:     dto.IngredientName,
		IngredientUnitName: dto.IngredientUnitName,
		Amount:             dto.Amount,
		OrderingIndex:      dto.OrderingIndex,
		Notes:              dto.Notes,
		CreatedAt:          dto.CreatedAt,
		CreatedBy:          dto.CreatedBy,
		UpdatedAt:          dto.UpdatedAt,
		UpdatedBy:          dto.UpdatedBy,
	}, nil
}

func (r *RecipeIngredientPostgresRepository) Delete(ctx context.Context, id uint64) error {
	dto, query := recipeIngredientDtoForDelete(id)

	_, err := r.db.NamedExecContext(ctx, query, &dto)
	if err == sql.ErrNoRows {
		return entity.ErrRecipeNotFound
	}

	return err
}

func recipeIngredientDtoForUpdate(id uint64, params usecase.RecipeIngredientParams, isDeleted *bool) (dto recipeIngredientDto, query string) {
	var qb strings.Builder

	qb.WriteString("UPDATE recipe_ingredients SET ")

	if params.Amount > 0 {
		qb.WriteString("amount = :amount, ")
		dto.Amount = params.Amount
	}

	if params.IngredientID != 0 {
		qb.WriteString("ingredient_id = :ingredient_id, ")
		dto.IngredientID = params.IngredientID
	}

	if params.IngredientName != "" {
		qb.WriteString("ingredient_name = :ingredient_name, ")
		dto.IngredientName = params.IngredientName
	}

	if params.IngredientUnitName != "" {
		qb.WriteString("ingredient_unit_name = :ingredient_unit_name, ")
		dto.IngredientUnitName = params.IngredientUnitName
	}

	if params.OrderingIndex > 0 {
		qb.WriteString("ordering_index = :ordering_index, ")
		dto.OrderingIndex = params.OrderingIndex
	}

	if params.Notes != "" {
		qb.WriteString("notes = :notes, ")
		dto.Notes = params.Notes
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

func recipeIngredientDtoForDelete(id uint64) (dto recipeIngredientDto, query string) {
	isDeleted := true
	return recipeIngredientDtoForUpdate(id, usecase.RecipeIngredientParams{}, &isDeleted)
}
