package models

import "time"

type CreateMovieRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Duration    int      `json:"duration" validate:"required,min=1"`
	Genres      string   `json:"genres" validate:"required"`
	WatchURL    string   `json:"watch_url" validate:"required,url"`
	Artists     []string `json:"artists" validate:"required"` // List of artist names
}

type Movie struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	Genres      string    `json:"genres"`
	WatchURL    string    `json:"watch_url"`
	Views       int       `json:"views"`
	Artists     []Artist  `json:"artists"` // Associated artists
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Artist struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
