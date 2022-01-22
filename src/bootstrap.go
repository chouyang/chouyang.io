package src

import (
	"chouyang.io/src/http"
	"chouyang.io/src/tools"
	"chouyang.io/src/types/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once

type Engine struct {
	*gin.Engine
	*gorm.DB
}

func (e *Engine) Ignite() *Engine {
	host := tools.Env("APP_HOST").(string)
	port := tools.Env("APP_PORT").(string)

	e.registerRoutes()

	fmt.Printf("Server started on %s:%s", host, port)
	_ = e.Run(":" + port)

	return e
}

func (e *Engine) registerRoutes() *Engine {
	routes := http.GetRoutes()
	for _, route := range routes {
		e.Handle(route.Method, route.Path, route.Handle)
	}

	return e
}

func (e *Engine) Migrate() {
	e.loadDB()
	_ = e.DB.AutoMigrate(&models.User{})
	_ = e.DB.AutoMigrate(&models.File{})
}

func (e *Engine) loadDB() {
	gin.DisableConsoleColor()
	var err error
	once.Do(func() {
		hostname := tools.Env("DB_HOST")
		portname := tools.Env("DB_PORT")
		username := tools.Env("DB_USER")
		password := tools.Env("DB_PASS")
		database := tools.Env("DB_NAME")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, portname, database)
		e.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})
}
