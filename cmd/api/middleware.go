package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) EnableCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Header("Access-Control-Allow-Origin", "https://supreme-halibut-v664446pgxqxhwxvr-3000.app.github.dev")
		// c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Origin", "https://go-react-frontend.vercel.app")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (app *Application) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, err := app.Auth.GetTokenAndVerify(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Next()
	}
}
