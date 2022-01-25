package handlers

import (
	"chouyang.io/src/errors"
	"chouyang.io/src/http/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

// FileHandler handles requests related to files
type FileHandler struct {
	Handler
	services.CodeService
}

// GetFileByPath returns the file at the given path if it exists
func (f *FileHandler) GetFileByPath(c *gin.Context) {
	root := fmt.Sprintf("%s%s", services.RealRoot, strings.ReplaceAll(c.Param("filepath"), "../", ""))
	root = strings.TrimRight(root, "/")

	fi, err := os.Stat(root)
	if err != nil {
		f.Error(c, errors.AccessDenied{})
		return
	}

	if fi.IsDir() {
		tree, err := f.ReadTree(root)
		if err != nil {
			f.Error(c, err)
			return
		}
		if tree == nil {
			f.Error(c, errors.AccessDenied{})
			return
		}

		c.JSON(200, tree)
	} else {
		file, err := f.ReadFile(root, true, fi)
		if err != nil {
			f.Error(c, err)
			return
		}

		c.JSON(200, file)
	}
}
