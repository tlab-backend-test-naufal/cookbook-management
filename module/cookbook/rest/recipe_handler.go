package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/guregu/null"

	"github.com/tlab-backend-test-naufal/cookbook-management/internal/libhttp"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/entity"
	"github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/internal/usecase"
)

type CreateRecipeRequest struct {
	RecipeRequest
	Ingredients RecipeIngredientRequests `json:"ingredients"`
}

type RecipeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  uint64 `json:"category_id"`
	Actor       string `json:"actor"`
}

type BulkCreateRecipeIngredientsRequest struct {
	RecipeID    uint64                   `json:"recipe_id"`
	Ingredients RecipeIngredientRequests `json:"ingredients"`
}

type RecipeIngredientRequests []RecipeIngredientRequest

type RecipeIngredientRequest struct {
	Amount             float64 `json:"amount"`
	IngredientID       uint64  `json:"ingredient_id"`
	IngredientName     string  `json:"ingredient_name"`
	IngredientUnitName string  `json:"ingredient_unit_name"`
	OrderingIndex      int     `json:"ordering_index"`
	Notes              string  `json:"notes"`
	Actor              string  `json:"actor"`
}

type RecipeResponse struct {
	ID          uint64      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CategoryID  uint64      `json:"category_id"`
	CreatedAt   time.Time   `json:"created_at"`
	CreatedBy   string      `json:"created_by"`
	UpdatedAt   null.Time   `json:"updated_at"`
	UpdatedBy   null.String `json:"updated_by"`
	IsDeleted   bool        `json:"is_deleted"`
}

type RecipeResponses struct {
	Data []RecipeResponse `json:"recipes"`
}

type RecipeIngredientResponse struct {
	ID                 uint64      `json:"id"`
	RecipeID           uint64      `json:"recipe_id,omitempty"`
	IngredientID       uint64      `json:"ingredient_id,omitempty"`
	IngredientName     string      `json:"ingredient_name"`
	IngredientUnitName string      `json:"ingredient_unit_name"`
	Amount             float64     `json:"amount"`
	OrderingIndex      int         `json:"ordering_index"`
	Notes              string      `json:"notes"`
	CreatedAt          time.Time   `json:"created_at"`
	CreatedBy          string      `json:"created_by"`
	UpdatedAt          null.Time   `json:"updated_at"`
	UpdatedBy          null.String `json:"updated_by"`
	IsDeleted          bool        `json:"is_deleted"`
}

type RecipeIngredientResponses struct {
	Data []RecipeIngredientResponse `json:"recipe_units"`
}

type GetSummaryResponse struct {
	RecipeResponse
	Ingredients []RecipeIngredientResponse `json:"ingredients"`
}

// CreateRecipe is a create recipe handler
func (h *CookbookHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var req CreateRecipeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	params := usecase.CreateRecipeParams{
		RecipeParams: usecase.RecipeParams{
			Name:        req.Name,
			Description: req.Description,
			CategoryID:  req.CategoryID,
			Actor:       req.Actor,
		},
	}

	for _, ingredient := range req.Ingredients {
		params.Ingredients = append(params.Ingredients, usecase.RecipeIngredientParams{
			Amount:             ingredient.Amount,
			IngredientID:       ingredient.IngredientID,
			IngredientName:     ingredient.IngredientName,
			IngredientUnitName: ingredient.IngredientUnitName,
			OrderingIndex:      ingredient.OrderingIndex,
			Notes:              ingredient.Notes,
			Actor:              req.Actor,
		})
	}

	err = h.recipeUsecase.CreateRecipe(r.Context(), params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithMessage(w, http.StatusOK, "successfully created recipe")
}

// UpdateRecipe is a update recipe handler.
func (h *CookbookHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	var req RecipeRequest
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("id cannot be empty"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	params := normalizeUpdateRecipeRequest(req)

	recipe, err := h.recipeUsecase.UpdateRecipe(r.Context(), id, params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, recipeResponseFromEntity(recipe))
	return
}

// DeleteRecipe is a delete recipe handler.
func (h *CookbookHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("id cannot be empty"))
		return
	}

	err = h.recipeUsecase.DeleteRecipe(r.Context(), id)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithMessage(w, http.StatusOK, "successfully deleted recipe")
	return
}

// BulkCreateRecipeIngredients is a create recipe recipe handler
func (h *CookbookHandler) BulkCreateRecipeIngredients(w http.ResponseWriter, r *http.Request) {
	var req BulkCreateRecipeIngredientsRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	var params usecase.BulkRecipeIngredientParams
	for _, i := range req.Ingredients {
		params = append(params, usecase.RecipeIngredientParams{
			Amount:             i.Amount,
			IngredientID:       i.IngredientID,
			IngredientName:     i.IngredientName,
			IngredientUnitName: i.IngredientUnitName,
			OrderingIndex:      i.OrderingIndex,
			Notes:              i.Notes,
			Actor:              i.Actor,
		})
	}

	err = h.recipeUsecase.BulkCreateRecipeIngredients(r.Context(), req.RecipeID, params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, "successfully created recipe ingredients")
}

// UpdateRecipeIngredient is a update recipe recipe handler.
func (h *CookbookHandler) UpdateRecipeIngredient(w http.ResponseWriter, r *http.Request) {
	var req RecipeIngredientRequest

	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("id cannot be empty"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	params := normalizeUpdateRecipeIngredientRequest(req)

	recipe, err := h.recipeUsecase.UpdateRecipeIngredient(r.Context(), id, params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, recipeIngredientResponseFromEntity(recipe))
	return
}

// DeleteRecipeIngredient is a delete recipe recipe handler.
func (h *CookbookHandler) DeleteRecipeIngredient(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("id cannot be empty"))
		return
	}

	err = h.recipeUsecase.DeleteRecipeIngredient(r.Context(), id)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithMessage(w, http.StatusOK, "successfully deleted recipe")
	return
}

// ListRecipes is a list recipe handler
func (h *CookbookHandler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	rawOfs := query.Get("offset")
	rawLim := query.Get("limit")

	categoryID := query.Get("category_id")
	ingredientID := query.Get("ingredient_id")

	ofs, _ := strconv.Atoi(rawOfs)
	lim, _ := strconv.Atoi(rawLim)

	catID, _ := strconv.ParseUint(categoryID, 10, 64)
	ingID, _ := strconv.ParseUint(ingredientID, 10, 64)

	filter := usecase.ListRecipesFiter{
		CategoryID:   catID,
		IngredientID: ingID,
	}

	recipeUnits, err := h.recipeUsecase.ListRecipes(r.Context(), filter, lim, ofs)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	var resp RecipeResponses
	for _, iu := range recipeUnits {
		resp.Data = append(resp.Data, recipeResponseFromEntity(iu))
	}

	libhttp.WithJSON(w, http.StatusOK, resp)
	return
}

// GetRecipeSummary is a get summary handler
func (h *CookbookHandler) GetRecipeSummary(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("id cannot be empty"))
		return
	}

	summary, err := h.recipeUsecase.GetRecipeSummary(r.Context(), id)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, getSummaryResponseFromEntity(summary))
	return
}

// normalizeUpdateRecipeRequest converts input to usecase params
func normalizeUpdateRecipeRequest(input RecipeRequest) usecase.RecipeParams {
	params := usecase.RecipeParams{
		Actor: input.Actor,
	}

	if input.Name != "" {
		params.Name = input.Name
	}

	return params
}

// normalizeUpdateRecipeIngredientRequest converts input to usecase params
func normalizeUpdateRecipeIngredientRequest(input RecipeIngredientRequest) usecase.RecipeIngredientParams {
	params := usecase.RecipeIngredientParams{
		Actor: input.Actor,
	}

	return params
}

// recipeResponseFromEntity converts recipe entity to response
func recipeResponseFromEntity(ent *entity.Recipe) RecipeResponse {
	return RecipeResponse{
		ID:          ent.ID,
		Name:        ent.Name,
		CategoryID:  ent.CategoryID,
		Description: ent.Description,
		CreatedAt:   ent.CreatedAt,
		CreatedBy:   ent.CreatedBy,
		UpdatedAt:   ent.UpdatedAt,
		UpdatedBy:   ent.UpdatedBy,
		IsDeleted:   ent.IsDeleted,
	}
}

// recipeIngredientResponseFromEntity converts recipe recipe entity to response
func recipeIngredientResponseFromEntity(ent *entity.RecipeIngredient) RecipeIngredientResponse {
	return RecipeIngredientResponse{
		ID:                 ent.ID,
		RecipeID:           ent.RecipeID,
		IngredientID:       ent.IngredientID,
		IngredientName:     ent.IngredientName,
		IngredientUnitName: ent.IngredientUnitName,
		Amount:             ent.Amount,
		OrderingIndex:      ent.OrderingIndex,
		Notes:              ent.Notes,
		CreatedAt:          ent.CreatedAt,
		CreatedBy:          ent.CreatedBy,
		UpdatedAt:          ent.UpdatedAt,
		UpdatedBy:          ent.UpdatedBy,
		IsDeleted:          ent.IsDeleted,
	}
}

// getSummaryResponseFromEntity converts summary entity to response
func getSummaryResponseFromEntity(ent entity.RecipeSummary) GetSummaryResponse {
	var ingredientResponses []RecipeIngredientResponse

	for _, ingredient := range ent.Ingredients {
		ingredientResponses = append(ingredientResponses, RecipeIngredientResponse{
			ID:                 ingredient.ID,
			RecipeID:           ingredient.RecipeID,
			IngredientID:       ingredient.IngredientID,
			IngredientName:     ingredient.IngredientName,
			IngredientUnitName: ingredient.IngredientUnitName,
			Amount:             ingredient.Amount,
			OrderingIndex:      ingredient.OrderingIndex,
			Notes:              ingredient.Notes,
			CreatedAt:          ingredient.CreatedAt,
			CreatedBy:          ingredient.CreatedBy,
			UpdatedAt:          ingredient.UpdatedAt,
			UpdatedBy:          ingredient.UpdatedBy,
			IsDeleted:          ingredient.IsDeleted,
		})
	}

	return GetSummaryResponse{
		RecipeResponse: RecipeResponse{
			ID:          ent.ID,
			Name:        ent.Name,
			Description: ent.Description,
			CategoryID:  ent.CategoryID,
			CreatedAt:   ent.CreatedAt,
			CreatedBy:   ent.CreatedBy,
			UpdatedAt:   ent.UpdatedAt,
			UpdatedBy:   ent.UpdatedBy,
			IsDeleted:   ent.IsDeleted,
		},
		Ingredients: ingredientResponses,
	}
}
