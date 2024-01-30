package main

import (
	"github.com/gin-gonic/gin"
	limits "github.com/gin-contrib/size"
)

func (app *Application) Routes() *gin.Engine {
	// create router
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(app.EnableCORS())
	r.Use(limits.RequestSizeLimiter(int64(1024*1024)))  // 1 MB

	r.GET("/", app.Home)
	r.POST("/authenticate", app.Authenticate)
	r.GET("/movies", app.AllMovies)

	return r
}
