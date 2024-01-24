package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/Snow-00/go-react-movies-backend/internal/models"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {
	var movies []*models.Movie
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT
				id, title, release_date, runtime, mpaa_rating, description, COALESCE(image, ''), created_at, updated_at
			  FROM movies ORDER BY title
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}