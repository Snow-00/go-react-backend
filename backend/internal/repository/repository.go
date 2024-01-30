package repository

import (
	"database/sql"

	"github.com/Snow-00/go-react-movies-backend/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllMovies() ([]*models.Movie, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}
