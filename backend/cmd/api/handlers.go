package main

import (
	"net/http"

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
	movies, err := app.DB.AllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, movies)
}
