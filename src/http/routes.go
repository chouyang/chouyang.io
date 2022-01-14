package http

import (
	"chouyang.io/src/http/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(e *gin.Engine) {
	e.GET("/~/Workspace/chouyang.io/*filepath", handlers.GetFileByPath)
}
