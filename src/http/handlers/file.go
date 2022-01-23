package handlers

import (
	"chouyang.io/src/errors"
	"chouyang.io/src/tools"
	"chouyang.io/src/types/models"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var cwd string
var rootPath string

var ignoreList []string

func init() {
	cwd, _ = os.Getwd()
	cwd = path.Dir(cwd)

	rootPath = tools.Env("APP_ROOT").(string)

	ignoreList = []string{".env", ".DS_Store", ".git/", ".idea/", ".vscode/", "node_modules/"}
}

type FileHandler struct {
	Handler
}

func (f *FileHandler) GetFileByPath(c *gin.Context) {
	root := fmt.Sprintf("%s%s", cwd, strings.ReplaceAll(c.Param("filepath"), "../", ""))
	root = strings.TrimRight(root, "/")

	fi, err := os.Stat(root)
	if err != nil {
		f.Error(c, errors.AccessDenied{})
		return
	}

	if fi.IsDir() {
		tree, err := f.readTree(root)
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
		file, err := f.readFile(root, fi)
		if err != nil {
			f.Error(c, err)
			return
		}

		c.JSON(200, file)
	}
}

func (f *FileHandler) readFile(path string, fi os.FileInfo) (*models.File, errors.Throwable) {
	of, err := os.Open(path)
	defer of.Close()
	if err != nil {
		return nil, errors.NotFound{}
	}

	if f.isForbidden(of.Name()) {
		return nil, errors.AccessDenied{}
	}

	content := make([]byte, fi.Size())
	_, _ = of.Read(content)

	fm := models.File{
		Name:       f.trimName(of.Name()),
		Path:       f.trimPath(path),
		Size:       fi.Size(),
		Mime:       f.GetFileMime(of.Name()),
		Hash:       f.Md5(content),
		Permission: fi.Mode().String(),
		Content:    string(content),
		CreatedBy:  0,
		CreatedAt:  time.Now(),
		UpdatedAt:  fi.ModTime(),
	}

	return &fm, nil
}

func (f *FileHandler) readTree(path string) (*models.Tree, errors.Throwable) {
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.AccessDenied{}
	}
	if f.isForbidden(path) {
		return nil, nil
	}
	tree := models.Tree{
		Trees: nil,
		Files: nil,
		Name:  f.trimName(path),
		Path:  f.trimPath(path),
	}
	for _, item := range items {
		if item.IsDir() {
			subTree, _ := f.readTree(fmt.Sprintf("%s/%s", path, item.Name()))
			if subTree != nil {
				tree.Trees = append(tree.Trees, subTree)
			}

			continue
		}

		file, _ := f.readFile(fmt.Sprintf("%s/%s", path, item.Name()), item)
		if file != nil {
			tree.Files = append(tree.Files, file)
		}
	}

	if tree.Trees == nil && tree.Files == nil {
		return nil, nil
	}

	return &tree, nil
}

func (f *FileHandler) GetFileMime(name string) string {
	list := map[string]string{
		"js":         "application/javascript",
		"css":        "text/css",
		"html":       "text/html",
		"json":       "application/json",
		"jpg":        "image/jpeg",
		"jpeg":       "image/jpeg",
		"png":        "image/png",
		"gif":        "image/gif",
		"svg":        "image/svg+xml",
		"ico":        "image/x-icon",
		"txt":        "text/plain",
		"md":         "text/markdown",
		"go":         "application/go",
		"gitignore":  "text/gitignore",
		"env":        "text/env",
		"sum":        "text/sum",
		"sh":         "application/shell",
		"Dockerfile": "text/dockerfile",
		"yml":        "text/yaml",
		"mod":        "text/module",
	}

	split := strings.Split(name, ".")

	if len(split) >= len(split)-1 {
		if mime, ok := list[split[len(split)-1]]; ok {
			return mime
		}
	}

	return "application/octet-stream"
}

func (f *FileHandler) Md5(content []byte) string {
	hash := md5.New()
	hash.Write(content)

	return strings.ToUpper(fmt.Sprintf("%x", hash.Sum(nil)))
}

func (f *FileHandler) trimPath(path string) string {
	p := strings.Trim(strings.ReplaceAll(path, cwd, ""), "/")

	if p == "" {
		return "/"
	}

	return fmt.Sprintf("%s/%s", rootPath, p)
}

func (f *FileHandler) trimName(name string) string {
	n := strings.Split(f.trimPath(name), "/")

	if n[len(n)-1] == "" {
		return "/"
	}

	return n[len(n)-1]
}

func (f *FileHandler) isForbidden(path string) bool {
	path = strings.Replace(path, cwd, "", 1)
	for _, item := range ignoreList {
		if strings.HasSuffix(item, "/") && strings.Contains(path, item) {
			return true
		}

		if strings.HasSuffix(path, item) {
			return true
		}
	}

	return false
}
