package http

import (
	"chouyang.io/src/http/handlers"
	"chouyang.io/src/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(e *gin.Engine) {
	file := handlers.FileHandler{}

	e.GET(fmt.Sprintf("%s/*filepath", tools.Env("APP_ROOT")), file.GetFileByPath)
}
