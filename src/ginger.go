package src

import (
	"chouyang.io/src/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

var Turbine *gin.Engine
var once sync.Once

func init() {
	gin.DisableConsoleColor()
	once.Do(func() {
		Turbine = gin.Default()
	})

	http.RegisterRoutes(Turbine)
}

func Ignite() {
	fmt.Println("Server started on http://127.0.0.1:8080")
	_ = Turbine.Run("127.0.0.1:8080")
}
