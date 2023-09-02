package main

import (
	"github.com/go-chi/chi"
	"github.com/subosito/gotenv"
	"github.com/tlab-backend-test-naufal/cookbook-management/internal/config"
	cookbookConfig "github.com/tlab-backend-test-naufal/cookbook-management/module/cookbook/config"
	"log"
	"net/http"
)

func main() {
	_ = gotenv.Load()

	db, err := config.BuildPostgres()
	if err != nil {
		log.Fatal(err.Error())
	}

	cookbookHandler := cookbookConfig.RegisterCookbookHandler(db)

	mux := chi.NewRouter()

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/recipes/{id}/summary", cookbookHandler.GetRecipeSummary)
		r.Get("/recipes", cookbookHandler.ListRecipes)
		r.Post("/recipes", cookbookHandler.CreateRecipe)
		r.Patch("/recipes/{id}", cookbookHandler.UpdateRecipe)
		r.Delete("/recipes/{id}", cookbookHandler.DeleteRecipe)
		r.Post("/recipe-ingredients", cookbookHandler.BulkCreateRecipeIngredients)
		r.Patch("/recipe-ingredients/{id}", cookbookHandler.UpdateRecipeIngredient)
		r.Delete("/recipe-ingredients/{id}", cookbookHandler.DeleteRecipeIngredient)

		r.Get("/categories", cookbookHandler.ListCategories)
		r.Post("/categories", cookbookHandler.CreateCategory)
		r.Patch("/categories/{id}", cookbookHandler.UpdateCategory)
		r.Delete("/categories/{id}", cookbookHandler.DeleteCategory)

		r.Get("/ingredients", cookbookHandler.ListIngredients)
		r.Post("/ingredients", cookbookHandler.CreateIngredient)
		r.Patch("/ingredients/{id}", cookbookHandler.UpdateIngredient)
		r.Delete("/ingredients/{id}", cookbookHandler.DeleteIngredient)

		r.Get("/ingredient-units", cookbookHandler.ListIngredientUnits)
		r.Post("/ingredient-units", cookbookHandler.CreateIngredientUnit)
		r.Patch("/ingredient-units/{id}", cookbookHandler.UpdateIngredientUnit)
		r.Delete("/ingredient-units/{id}", cookbookHandler.DeleteIngredientUnit)
	})

	port := config.RestPort()

	log.Printf("Rest API is serving at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
