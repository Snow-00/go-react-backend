package main

import "github.com/gin-gonic/gin"

func (app *Application) Routes() *gin.Engine {
	// create router
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(app.EnableCORS())

	r.GET("/", app.Home)
	r.GET("/authenticate", app.Authenticate)
	r.GET("/movies", app.AllMovies)

	return r
}
