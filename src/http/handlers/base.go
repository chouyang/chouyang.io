package handlers

import (
	"chouyang.io/src/errors"
	"github.com/gin-gonic/gin"
)

func Okay(c *gin.Context, payload interface{}) {
	c.JSON(200, gin.H{
		"code":    errors.ApiOK,
		"message": "ok",
		"payload": payload,
	})
}

func Error(c *gin.Context, err errors.Throwable) {
	c.JSON(200, gin.H{
		"code":    err.GetCode(),
		"message": err.Error(),
	})
}
