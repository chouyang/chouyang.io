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

var Turbine *gin.Engine
var OilTank *gorm.DB
var once sync.Once
var err error

func init() {
	gin.DisableConsoleColor()
	once.Do(func() {
		hostname := tools.Env("DB_HOST")
		portname := tools.Env("DB_PORT")
		username := tools.Env("DB_USER")
		password := tools.Env("DB_PASS")
		database := tools.Env("DB_NAME")

		Turbine = gin.Default()
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, portname, database)
		OilTank, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})

	http.RegisterRoutes(Turbine)
}

func Ignite() {
	AutoMigrate()

	host := tools.Env("APP_HOST").(string)
	port := tools.Env("APP_PORT").(string)
	fmt.Printf("Server started on %s:%s", host, port)
	_ = Turbine.Run(":" + port)
}

func AutoMigrate() {
	_ = OilTank.AutoMigrate(&models.User{})
	_ = OilTank.AutoMigrate(&models.File{})
}
