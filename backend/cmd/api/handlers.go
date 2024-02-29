package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Snow-00/go-react-movies-backend/internal/graph"
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
	movie, err = app.GetPoster(movie)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	newID, err := app.DB.InsertMovie(movie)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	// now handle genres
	err = app.DB.UpdateMovieGenres(newID, movie.GenresArray)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "movie updated"})
}

func (app *Application) GetPoster(movie models.Movie) (models.Movie, error) {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalPages int `json:"total_pages"`
	}

	client := &http.Client{}
	apiUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", app.APIKey)

	req, err := http.NewRequest("GET", apiUrl+"&query="+url.QueryEscape(movie.Title), nil)
	if err != nil {
		return movie, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return movie, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Println(err)
		return movie, err
	}

	var responseObject TheMovieDB

	err = json.Unmarshal(bodyBytes, &responseObject)
	if err != nil {
		return movie, err
	}

	if len(responseObject.Results) > 0 {
		movie.Image = responseObject.Results[0].PosterPath
	}

	return movie, nil
}

func (app *Application) UpdateMovie(c *gin.Context) {
	var payload models.Movie

	err := c.BindJSON(&payload)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	movie, err := app.DB.OneMovie(payload.ID)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	// this method quite roundabout
	movie.Title = payload.Title
	movie.ReleaseDate = payload.ReleaseDate
	movie.Description = payload.Description
	movie.MPAARating = payload.MPAARating
	movie.RunTime = payload.RunTime
	movie.UpdatedAt = time.Now()

	err = app.DB.UpdateMovie(*movie)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	err = app.DB.UpdateMovieGenres(movie.ID, payload.GenresArray)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "movie updated"})
}

func (app *Application) DeleteMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	err = app.DB.DeleteMovie(id)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "movie deleted"})
}

func (app *Application) AllMoviesByGenre(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	movies, err := app.DB.AllMovies(id)
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (app *Application) MoviesGraphQL(c *gin.Context) {
	// we need to populate graph type w/ movies
	movies, _ := app.DB.AllMovies()

	// get the query from the request
	q, _ := io.ReadAll(c.Request.Body)
	query := string(q)

	// crete a new var of type *graph.Graph
	g := graph.New(movies)

	// set the query string on the var
	g.QueryString = query

	// performs the query
	resp, err := g.Query()
	if err != nil {
		app.ErrorJSON(c, err)
		return
	}

	// send the response
	j, _ := json.MarshalIndent(resp, "", "\t")
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(j)
}
