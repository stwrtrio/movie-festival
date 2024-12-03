package models

import "time"

type CreateMovieRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Duration    int      `json:"duration" validate:"required,min=1"`
	Genres      []string `json:"genres" validate:"min=1,dive,required"`
	WatchURL    string   `json:"watch_url" validate:"required,url"`
	Artists     []string `json:"artists" validate:"min=1,dive,required"` // List of artist names
}

type Movie struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	Genres      []Genre   `json:"genres"`
	WatchURL    string    `json:"watch_url"`
	Views       int       `json:"views"`
	Artists     []Artist  `json:"artists"` // Associated artists
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GenreView struct {
	Name      string `json:"name"`
	ViewCount int64  `json:"total_views"`
}

type MovieView struct {
	MovieID      int       `json:"movie_id"`
	UserID       int       `json:"user_id,omitempty"`
	ViewCount    int       `json:"view_count"`
	LastViewedAt time.Time `json:"last_viewed_at"`
}

type Vote struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	MovieID   string    `json:"movie_id"`
	CreatedAt time.Time `json:"created_at"`
}
