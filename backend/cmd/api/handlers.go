package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Snow-00/go-react-movies-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		app.ErrorJSON(c, errors.New("invalid credentials (email)"), http.StatusForbidden)
		return
	}

	//check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJSON(c, errors.New("invalid credentials (pass)"), http.StatusForbidden)
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
	app.Auth.RefreshCookie(c, tokens.RefreshToken)

	c.JSON(http.StatusAccepted, tokens)
}

func (app *Application) RefreshToken(c *gin.Context) {
	// looping through cookies that were sent to backend
	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == app.Auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.ErrorJSON(c, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get user id from refresh token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.ErrorJSON(c, errors.New("unknowm user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.ErrorJSON(c, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := JWTUser{
				ID:        userID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.Auth.GenerateTokenPair(&u)
			if err != nil {
				app.ErrorJSON(c, errors.New("error generating tokens"))
				return
			}

			app.Auth.RefreshCookie(c, tokenPairs.RefreshToken)

			c.JSON(http.StatusOK, tokenPairs)
		}
	}
}

func (app *Application) Logout(c *gin.Context) {
	app.Auth.ExpiredRefreshCookie(c)
	c.Writer.WriteHeader(http.StatusAccepted)
}

func (app *Application) MovieCatalog(c *gin.Context) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (app *Application) GetMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (app *Application) MovieForEdit(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	payload := struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		movie,
		genres,
	}

	c.JSON(http.StatusOK, payload)
}

func (app *Application) AllGenres(c *gin.Context) {
	genres, err := app.DB.AllGenres()
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, genres)
}

func (app *Application) InsertMovie(c *gin.Context) {
	var movie models.Movie

	err := c.BindJSON(&movie)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	// try to get image

	// now handle genres

	c.JSON(http.StatusAccepted, gin.H{"message": "movie updated"})
}
