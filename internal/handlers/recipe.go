package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Fitness-Service/recipe-service/internal/models"
	"github.com/Fitness-Service/recipe-service/internal/repository"
)

type RecipeHandler struct {
	repo *repository.RecipeRepo
}

func NewRecipeHandler(repo *repository.RecipeRepo) *RecipeHandler {
	return &RecipeHandler{repo: repo}
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&recipe); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) SuggestByIngredients(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Ingredients []string `json:"ingredients"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	recipes, err := h.repo.SuggestByIngredients(req.Ingredients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) SuggestByGoal(w http.ResponseWriter, r *http.Request) {
	goal := r.URL.Query().Get("goal")
	if goal == "" {
		http.Error(w, "goal required", http.StatusBadRequest)
		return
	}

	recipes, err := h.repo.SuggestByGoal(goal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipes)
}
