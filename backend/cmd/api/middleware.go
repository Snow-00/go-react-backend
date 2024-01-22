package main

import "github.com/gin-gonic/gin"

func (app *Application) EnableCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://*")

		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Method", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")

			// THIS IS JUST IF SOMETHING GOES WRONG
			// return
		} // else { c.Next() }

		c.Next()
	}
}
