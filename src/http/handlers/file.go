package handlers

import (
	"chouyang.io/src/errors"
	"chouyang.io/src/types/models"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var cwd string

func init() {
	cwd, _ = os.Getwd()
	cwd = fmt.Sprintf("%s/", cwd)
}

type FileHandler struct {
	Handler
}

func (f *FileHandler) GetFileByPath(c *gin.Context) {
	path := fmt.Sprintf("%s%s", cwd, strings.ReplaceAll(c.Param("filepath"), "../", ""))

	fi, err := os.Stat(path)
	if err != nil {
		f.Error(c, errors.NotFound{})
		return
	}

	if fi.IsDir() {
		tree, err := f.readTree(path)
		if err != nil {
			f.Error(c, err)
			return
		}
		c.JSON(200, tree)
	} else {
		file, err := f.readFile(path, fi)
		if err != nil {
			f.Error(c, err)
			return
		}

		c.Data(200, file.Mime, []byte(file.Content))
	}
}

func (f *FileHandler) readFile(path string, fi os.FileInfo) (*models.File, errors.Throwable) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, errors.NotFound{}
	}

	content := make([]byte, fi.Size())
	_, _ = file.Read(content)

	return &models.File{
		ID:         0,
		Name:       file.Name(),
		Path:       strings.Replace(path, cwd, "", 1),
		Size:       fi.Size(),
		Mime:       f.GetFileMime(file.Name()),
		Hash:       f.Md5(content),
		Permission: fi.Mode().String(),
		Content:    string(content),
		CreatedBy:  0,
		CreatedAt:  time.Now(),
		UpdatedAt:  fi.ModTime(),
	}, nil
}

func (f *FileHandler) readTree(path string) (*models.Tree, errors.Throwable) {
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.AccessDenied{}
	}
	tree := &models.Tree{
		Trees: nil,
		Files: nil,
		Name:  path,
		Path:  strings.Replace(path, cwd, "", 1),
	}
	for _, item := range items {
		switch item.Name() {
		case ".", "..", ".git", ".env", ".idea", ".vscode", ".DS_Store", "node_modules":
			continue
		}
		if item.IsDir() {
			subTree, err := f.readTree(fmt.Sprintf("%s/%s", path, item.Name()))
			if err != nil {
				return nil, err
			}
			tree.Trees = append(tree.Trees, subTree)
		} else {
			file, err := f.readFile(fmt.Sprintf("%s/%s", path, item.Name()), item)
			if err != nil {
				return nil, err
			}
			tree.Files = append(tree.Files, file)
		}
	}

	return tree, nil
}

func (f *FileHandler) GetFileMime(name string) string {
	list := map[string]string{
		"js":   "application/javascript",
		"css":  "text/css",
		"html": "text/html",
		"json": "application/json",
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"svg":  "image/svg+xml",
		"ico":  "image/x-icon",
		"txt":  "text/plain",
		"md":   "text/markdown",
		"go":   "text/go",
	}

	ext := strings.ToLower(strings.TrimLeft(name, "."))

	if mime, ok := list[ext]; ok {
		return mime
	}

	return "application/octet-stream"
}

func (f *FileHandler) Md5(content []byte) string {
	hash := md5.New()
	hash.Write(content)

	return fmt.Sprintf("%x", hash.Sum(nil))
}
