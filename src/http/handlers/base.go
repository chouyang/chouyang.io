package handlers

import (
	"chouyang.io/src/errors"
	"github.com/gin-gonic/gin"
)

// Handler implements common handler functions.
type Handler struct {
}

// Okay is a helper function to return an ApiOK response.
func (h *Handler) Okay(c *gin.Context, payload interface{}) {
	c.JSON(200, gin.H{
		"code":    errors.ApiOK,
		"message": "ok",
		"payload": payload,
	})
}

// Error is a helper function to return a Throwable error response.
func (h *Handler) Error(c *gin.Context, err errors.Throwable) {
	c.JSON(200, gin.H{
		"code":    err.GetCode(),
		"message": err.Error(),
	})
}
