package models

import "time"

type Recipe struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Ingredients []string  `json:"ingredients" db:"ingredients"`
	Calories    int       `json:"calories" db:"calories"`
	Protein     float64   `json:"protein" db:"protein"`
	Fats        float64   `json:"fats" db:"fats"`
	Carbs       float64   `json:"carbs" db:"carbs"`
	Goal        string    `json:"goal" db:"goal"`
	Tags        []string  `json:"tags" db:"tags"`
	PrepTime    int       `json:"prep_time" db:"prep_time"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
