package main

import (
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func (app *Application) Routes() *gin.Engine {
	// create router
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(app.EnableCORS())
	r.Use(limits.RequestSizeLimiter(int64(1024 * 1024))) // 1 MB

	r.GET("/", app.Home)
	r.POST("/authenticate", app.Authenticate)
	r.GET("/refresh", app.RefreshToken)
	r.GET("/logout", app.Logout)
	r.GET("/movies", app.AllMovies)
	r.GET("/movies/:id", app.GetMovie)

	r.GET("/genres", app.AllGenres)

	authorized := r.Group("/admin", app.AuthRequired())
	{
		authorized.GET("/movies", app.MovieCatalog)
		authorized.GET("/movies/:id", app.MovieForEdit)
	}

	return r
}
