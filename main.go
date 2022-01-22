package main

import (
	"chouyang.io/src"
	"github.com/gin-gonic/gin"
)

func main() {
	turbine := &src.Engine{Engine: gin.Default()}
	turbine.Migrate()

	turbine.Ignite()
}
