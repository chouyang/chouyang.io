package http

import (
	"chouyang.io/src/http/handlers"
	"chouyang.io/src/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Method string
	Path   string
	Handle gin.HandlerFunc
}

func GetRoutes() []*Route {
	file := handlers.FileHandler{}
	var routes []*Route
	routes = append(routes, loadGet(fmt.Sprintf("%s/*filepath", tools.Env("APP_ROOT")), file.GetFileByPath))

	return routes
}

func loadGet(path string, handle gin.HandlerFunc) *Route {
	return request("GET", path, handle)
}

func loadPost(path string, handle gin.HandlerFunc) *Route {
	return request("POST", path, handle)
}

func loadPut(path string, handle gin.HandlerFunc) *Route {
	return request("PUT", path, handle)
}

func loadDelete(path string, handle gin.HandlerFunc) *Route {
	return request("DELETE", path, handle)
}

func request(method string, path string, handle gin.HandlerFunc) *Route {
	return &Route{
		Method: method,
		Path:   path,
		Handle: handle,
	}
}
