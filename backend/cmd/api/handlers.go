package main

import (
	"errors"
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
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.BindJSON(&requestPayload)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	// validate user against db (doesnt know the user)
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.ErrorJSON(c, errors.New("invalid credentials"), http.StatusForbidden)
		return
	}

	//check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJSON(c, errors.New("invalid credentials"), http.StatusForbidden)
		return
	}

	// create jwt user
	u := JWTUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.Auth.GenerateTokenPair(&u)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	// set refresh cookie
	var j *Auth

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		j.CookieName,
		tokens.RefreshToken,
		int(j.RefreshExpiry.Seconds()),
		j.CookiePath,
		j.CookieDomain,
		true,
		true,
	)

	c.JSON(http.StatusAccepted, tokens)
}
