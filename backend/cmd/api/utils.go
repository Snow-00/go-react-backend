package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

func (app *Application) ErrorJSON(c *gin.Context, err error, status ...int) {
	statusCode := http.StatusBadRequest

	// if there is status in param, use this
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload ErrorMessage
	payload.Error = true
	payload.Message = err.Error()

	c.JSON(statusCode, payload)
}

func (j *Auth) RefreshCookie(c *gin.Context, refreshToken string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		j.CookieName,
		refreshToken,
		int(j.RefreshExpiry.Seconds()),
		j.CookiePath,
		j.CookieDomain,
		true,
		true,
	)
}