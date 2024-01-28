package main

import (
	"log"
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

	c.JSON(http.StatusOK, payload)
}

func (app *Application) AllMovies(c *gin.Context) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (app *Application) Authenticate(c *gin.Context) {
	// read json payload

	// validate user against db

	//check password

	// create jwt user
	u := JWTUser{
		ID:        1,
		FirstName: "Admin",
		LastName:  "User",
	}

	// generate tokens
	tokens, err := app.Auth.GenerateTokenPair(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log.Println(tokens.Token)
	c.JSON(http.StatusOK, []byte(tokens.Token))
}
