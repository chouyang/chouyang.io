package http

import (
	"chouyang.io/src/http/handlers"
	"chouyang.io/src/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Route defines an http route.
type Route struct {
	Method string
	Path   string
	Handle gin.HandlerFunc
}

// GetRoutes returns a list of routes.
// we define our routes here.
func GetRoutes() []*Route {
	file := handlers.FileHandler{}

	return []*Route{
		request("GET", fmt.Sprintf("%s/*filepath", tools.Env("APP_ROOT")), file.GetFileByPath),
	}
}

// request instantiates a route.
func request(method string, path string, handle gin.HandlerFunc) *Route {
	return &Route{
		Method: method,
		Path:   path,
		Handle: handle,
	}
}
