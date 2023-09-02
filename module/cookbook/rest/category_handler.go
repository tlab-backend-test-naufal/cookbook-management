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

type CategoryRequest struct {
	Name  string `json:"name"`
	Actor string `json:"actor"`
}

type CategoryResponse struct {
	ID        uint64      `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	UpdatedAt null.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
	IsDeleted bool        `json:"is_deleted"`
}

type CategoriesResponse struct {
	CategoryResponse []CategoryResponse `json:"categories"`
}

// CreateCategory is a create category handler
func (h *CookbookHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	params := usecase.CategoryParams{
		Name:  req.Name,
		Actor: req.Actor,
	}

	category, err := h.categoryUsecase.CreateCategory(r.Context(), params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, categoryResponseFromEntity(category))
}

// UpdateCategory is a update category handler.
func (h *CookbookHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest
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

	params := normalizeUpdateCategoryRequest(req)

	category, err := h.categoryUsecase.UpdateCategory(r.Context(), id, params)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithJSON(w, http.StatusOK, categoryResponseFromEntity(category))
	return
}

// DeleteCategory is a delete category handler.
func (h *CookbookHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		libhttp.WithError(w, http.StatusBadRequest, fmt.Errorf("id cannot be empty"))
		return
	}

	err = h.categoryUsecase.DeleteCategory(r.Context(), id)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	libhttp.WithMessage(w, http.StatusOK, "successfully deleted category")
	return
}

// ListCategories is a list category handler
func (h *CookbookHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	rawOfs := query.Get("offset")
	rawLim := query.Get("limit")

	ofs, _ := strconv.Atoi(rawOfs)
	lim, _ := strconv.Atoi(rawLim)

	categories, err := h.categoryUsecase.ListCategories(r.Context(), lim, ofs)
	if err != nil {
		libhttp.WithError(w, http.StatusInternalServerError, err)
		return
	}

	var resp CategoriesResponse
	for _, c := range categories {
		resp.CategoryResponse = append(resp.CategoryResponse, categoryResponseFromEntity(c))
	}

	libhttp.WithJSON(w, http.StatusOK, resp)
	return
}

// normalizeUpdateCategoryRequest converts input to usecase params
func normalizeUpdateCategoryRequest(input CategoryRequest) usecase.CategoryParams {
	params := usecase.CategoryParams{
		Actor: input.Actor,
	}

	if input.Name != "" {
		params.Name = input.Name
	}

	return params
}

// categoryResponseFromEntity converts category entity to response
func categoryResponseFromEntity(ent *entity.Category) CategoryResponse {
	return CategoryResponse{
		ID:        ent.ID,
		Name:      ent.Name,
		CreatedAt: ent.CreatedAt,
		CreatedBy: ent.CreatedBy,
		UpdatedAt: ent.UpdatedAt,
		UpdatedBy: ent.UpdatedBy,
		IsDeleted: ent.IsDeleted,
	}
}
