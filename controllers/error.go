package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorController generates pages with a 404 response
func ErrorController(c *gin.Context) {
	renderSPA(c, http.StatusNotFound)
}
