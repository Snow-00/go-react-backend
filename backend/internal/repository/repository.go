package repository

import "github.com/Snow-00/go-react-movies-backend/internal/models"

type DatabaseRepo interface {
	AllMovies() ([]*models.Movie, error)
}