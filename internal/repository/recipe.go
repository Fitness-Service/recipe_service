package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Fitness-Service/recipe-service/internal/models"
)

type RecipeRepo struct {
	db *sqlx.DB
}

func NewRecipeRepo(connString string) *RecipeRepo {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal("DB connect failed:", err)
	}

	// миграция
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS recipes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		ingredients TEXT[],
		calories INT,
		protein FLOAT,
		fats FLOAT,
		carbs FLOAT,
		goal TEXT,
		tags TEXT[],
		prep_time INT,
		created_at TIMESTAMP DEFAULT NOW()
	)`)

	return &RecipeRepo{db: db}
}

func (r *RecipeRepo) Close() {
	r.db.Close()
}

func (r *RecipeRepo) Create(recipe *models.Recipe) error {
	return r.db.QueryRowx(`
		INSERT INTO recipes (name, ingredients, calories, protein, fats, carbs, goal, tags, prep_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`,
		recipe.Name, recipe.Ingredients, recipe.Calories,
		recipe.Protein, recipe.Fats, recipe.Carbs,
		recipe.Goal, recipe.Tags, recipe.PrepTime,
	).Scan(&recipe.ID, &recipe.CreatedAt)
}

func (r *RecipeRepo) GetAll() ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.db.Select(&recipes, "SELECT * FROM recipes ORDER BY created_at DESC")
	return recipes, err
}

func (r *RecipeRepo) SuggestByIngredients(ingredients []string) ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.db.Select(&recipes,
		"SELECT * FROM recipes WHERE ingredients && $1",
		ingredients,
	)
	return recipes, err
}

func (r *RecipeRepo) SuggestByGoal(goal string) ([]models.Recipe, error) {
	var recipes []models.Recipe
	err := r.db.Select(&recipes,
		"SELECT * FROM recipes WHERE goal = $1",
		goal,
	)
	return recipes, err
}
