package main

import (
	"time"

	"github.com/Snow-00/go-react-movies-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (app *Application) Home(c *gin.Context) {
	// this is only for testing
	payload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}

	c.JSON(200, payload)
}

func (app *Application) AllMovies(c *gin.Context) {
	var movies []models.Movie

	rdH, _ := time.Parse("2006-01-02", "1986-03-07")
	rdR, _ := time.Parse("2006-01-02", "1981-06-12")

	// this is just for testing
	highlander := models.Movie{
		ID:          1,
		Title:       "Highlander",
		ReleaseDate: rdH,
		MPAARating:  "R",
		RunTime:     116, // minutes
		Description: "a very nice movie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	movies = append(movies, highlander)

	rotla := models.Movie{
		ID:          2,
		Title:       "Raiders of the Lost Ark",
		ReleaseDate: rdR,
		MPAARating:  "PG-13",
		RunTime:     115, // minutes
		Description: "another very nice movie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	movies = append(movies, rotla)

	c.JSON(200, movies)
}
