package services

import (
	"chouyang.io/src/errors"
	"chouyang.io/src/tools"
	"chouyang.io/src/types/models"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// CodeService handles requests related to files
type CodeService struct {
}

// ReadFile reads the file at the given path and returns a File object
func (cs *CodeService) ReadFile(path string, fi os.FileInfo) (*models.File, errors.Throwable) {
	of, err := os.Open(path)
	defer of.Close()
	if err != nil {
		return nil, errors.NotFound{}
	}

	if cs.IsForbidden(of.Name()) {
		return nil, errors.AccessDenied{}
	}

	content := make([]byte, fi.Size())
	_, _ = of.Read(content)

	fm := models.File{
		Name:       cs.TrimName(of.Name()),
		Path:       cs.TrimPath(path),
		Size:       fi.Size(),
		Mime:       cs.GetFileMime(of.Name()),
		Hash:       tools.Md5(content),
		Permission: fi.Mode().String(),
		Content:    string(content),
		CreatedBy:  0,
		CreatedAt:  time.Now(),
		UpdatedAt:  fi.ModTime(),
	}

	return &fm, nil
}

// ReadTree recursively reads the directory at the given path and returns a Tree object
func (cs *CodeService) ReadTree(path string) (*models.Tree, errors.Throwable) {
	if cs.IsForbidden(path) {
		return nil, nil
	}
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.AccessDenied{}
	}
	tree := models.Tree{
		Trees:     []*models.Tree{},
		Files:     []string{},
		Name:      cs.TrimName(path),
		Path:      cs.TrimPath(path),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	for _, item := range items {
		if item.IsDir() {
			subTree, _ := cs.ReadTree(fmt.Sprintf("%s/%s", path, item.Name()))
			if subTree != nil {
				tree.Trees = append(tree.Trees, subTree)
			}

			continue
		}

		if cs.IsForbidden(fmt.Sprintf("%s/%s", path, item.Name())) {
			continue
		}

		tree.Files = append(tree.Files, item.Name())
	}

	if tree.Trees == nil && tree.Files == nil {
		return nil, nil
	}

	return &tree, nil
}

// GetFileMime returns the mime type of the file based on its extension
func (cs *CodeService) GetFileMime(name string) string {
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

// TrimPath replaces the root path with the virtual root path
func (cs *CodeService) TrimPath(path string) string {
	p := strings.Trim(strings.ReplaceAll(path, RealRoot, ""), "/")

	if p == "" {
		return "/"
	}

	return fmt.Sprintf("%s/%s", VirtualRoot, p)
}

// TrimName trims out the root path from the file name
func (cs *CodeService) TrimName(name string) string {
	n := strings.Split(cs.TrimPath(name), "/")

	if n[len(n)-1] == "" {
		return "chouyang.io"
	}

	return n[len(n)-1]
}

// IsForbidden checks if the file is in the ignoreList
func (cs *CodeService) IsForbidden(path string) bool {
	path = strings.Replace(path, RealRoot, "", 1)
	for _, item := range ignoreList {
		if strings.HasSuffix(item, "/") {
			// check if segments of the path match the ignoreList
			segments := strings.Split(path, "/")
			for _, segment := range segments {
				if strings.Trim(segment, "/") == strings.Trim(item, "/") {
					return true
				}
			}
		}

		if strings.HasSuffix(path, item) {
			return true
		}
	}

	return false
}
