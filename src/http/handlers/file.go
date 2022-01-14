package handlers

import (
	"chouyang.io/src/errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func GetFileByPath(c *gin.Context) {
	cwd, _ := os.Getwd()
	filepath := fmt.Sprintf("%s%s", cwd, strings.ReplaceAll(c.Param("filepath"), "../", ""))

	content, err := readFile(filepath)
	if err != nil {
		Error(c, err)
		return
	}

	Okay(c, content)
}

func readFile(filepath string) (string, errors.Throwable) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return "", errors.NotFound{}
	}

	if fi.IsDir() {
		return "", errors.AccessDenied{}
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", errors.NotFound{}
	}

	content := make([]byte, fi.Size())
	_, _ = file.Read(content)

	return string(content), nil
}
