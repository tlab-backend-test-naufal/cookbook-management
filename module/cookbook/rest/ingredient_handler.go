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

type IngredientRequest struct {
	Name  string `json:"name"`
	Actor string `json:"actor"`
}

type IngredientResponse struct {
	ID        uint64      `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	UpdatedAt null.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
	IsDeleted bool        `json:"is_deleted"`
}

type IngredientsResponse struct {
	Data []IngredientResponse `json:"ingredients"`
}

type IngredientUnitRequest struct {
	Name  string `json:"name"`
	Actor string `json:"actor"`
}

type IngredientUnitResponse struct {
	ID        uint64      `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	UpdatedAt null.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
	IsDeleted bool        `json:"is_deleted"`
}

type IngredientUnitResponses struct {
	Data []IngredientUnitResponse `json:"ingredient_units"`
}

// CreateIngredient is a create ingredient handler
func (h *CookbookHandler) CreateIngredient(w http.ResponseWriter, r *http.Request) {
	var req IngredientRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	params := usecase.IngredientParams{
		Name:  req.Name,
		Actor: req.Actor,
	}

	ingredient, err := h.ingredientUsecase.CreateIngredient(r.Context(), params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, ingredientResponseFromEntity(ingredient))
}

// UpdateIngredient is a update ingredient handler.
func (h *CookbookHandler) UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	var req IngredientRequest

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

	params := normalizeUpdateIngredientRequest(req)

	ingredient, err := h.ingredientUsecase.UpdateIngredient(r.Context(), id, params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, ingredientResponseFromEntity(ingredient))
	return
}

// DeleteIngredient is a delete ingredient handler.
func (h *CookbookHandler) DeleteIngredient(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("id cannot be empty"))
		return
	}

	err = h.ingredientUsecase.DeleteIngredient(r.Context(), id)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithMessage(w, http.StatusOK, "successfully deleted ingredient unit")
	return
}

// ListIngredients is a list ingredient handler
func (h *CookbookHandler) ListIngredients(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	rawOfs := query.Get("offset")
	rawLim := query.Get("limit")

	ofs, _ := strconv.Atoi(rawOfs)
	lim, _ := strconv.Atoi(rawLim)

	ingredients, err := h.ingredientUsecase.ListIngredients(r.Context(), lim, ofs)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	var resp IngredientsResponse
	for _, i := range ingredients {
		resp.Data = append(resp.Data, ingredientResponseFromEntity(i))
	}

	libhttp.WithJSON(w, http.StatusOK, resp)
	return
}

// CreateIngredientUnit is a create ingredient handler
func (h *CookbookHandler) CreateIngredientUnit(w http.ResponseWriter, r *http.Request) {
	var req IngredientUnitRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	params := usecase.IngredientUnitParams{
		Name:  req.Name,
		Actor: req.Actor,
	}

	ingredient, err := h.ingredientUsecase.CreateIngredientUnit(r.Context(), params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, ingredientUnitResponseFromEntity(ingredient))
}

// UpdateIngredientUnit is a update ingredient handler.
func (h *CookbookHandler) UpdateIngredientUnit(w http.ResponseWriter, r *http.Request) {
	var req IngredientUnitRequest

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

	params := normalizeUpdateIngredientUnitRequest(req)

	ingredient, err := h.ingredientUsecase.UpdateIngredientUnit(r.Context(), id, params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, ingredientUnitResponseFromEntity(ingredient))
	return
}

// DeleteIngredientUnit is a delete ingredient handler.
func (h *CookbookHandler) DeleteIngredientUnit(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("id cannot be empty"))
		return
	}

	err = h.ingredientUsecase.DeleteIngredientUnit(r.Context(), id)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithMessage(w, http.StatusOK, "successfully deleted ingredient")
	return
}

// ListIngredientUnits is a list ingredient handler
func (h *CookbookHandler) ListIngredientUnits(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	rawOfs := query.Get("offset")
	rawLim := query.Get("limit")

	ofs, _ := strconv.Atoi(rawOfs)
	lim, _ := strconv.Atoi(rawLim)

	ingredientUnits, err := h.ingredientUsecase.ListIngredientUnits(r.Context(), lim, ofs)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	var resp IngredientUnitResponses
	for _, iu := range ingredientUnits {
		resp.Data = append(resp.Data, ingredientUnitResponseFromEntity(iu))
	}

	libhttp.WithJSON(w, http.StatusOK, resp)
	return
}

// normalizeUpdateIngredientRequest converts input to usecase params
func normalizeUpdateIngredientRequest(input IngredientRequest) usecase.IngredientParams {
	params := usecase.IngredientParams{
		Actor: input.Actor,
	}

	if input.Name != "" {
		params.Name = input.Name
	}

	return params
}

// normalizeUpdateIngredientUnitRequest converts input to usecase params
func normalizeUpdateIngredientUnitRequest(input IngredientUnitRequest) usecase.IngredientUnitParams {
	params := usecase.IngredientUnitParams{
		Actor: input.Actor,
	}

	if input.Name != "" {
		params.Name = input.Name
	}

	return params
}

// ingredientResponseFromEntity converts ingredient entity to response
func ingredientResponseFromEntity(ent *entity.Ingredient) IngredientResponse {
	return IngredientResponse{
		ID:        ent.ID,
		Name:      ent.Name,
		CreatedAt: ent.CreatedAt,
		CreatedBy: ent.CreatedBy,
		UpdatedAt: ent.UpdatedAt,
		UpdatedBy: ent.UpdatedBy,
		IsDeleted: ent.IsDeleted,
	}
}

// ingredientUnitResponseFromEntity converts ingredient unit entity to response
func ingredientUnitResponseFromEntity(ent *entity.IngredientUnit) IngredientUnitResponse {
	return IngredientUnitResponse{
		ID:        ent.ID,
		Name:      ent.Name,
		CreatedAt: ent.CreatedAt,
		CreatedBy: ent.CreatedBy,
		UpdatedAt: ent.UpdatedAt,
		UpdatedBy: ent.UpdatedBy,
		IsDeleted: ent.IsDeleted,
	}
}
