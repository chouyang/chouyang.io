package services

import (
	"chouyang.io/src/errors"
	"chouyang.io/src/tools"
	"chouyang.io/src/types/models"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var projectRoot = "chouyang.io"

// CodeService handles requests related to files
type CodeService struct {
}

// ReadFile reads the file at the given path and returns a File object
func (cs *CodeService) ReadFile(path string, loadContent bool, fi os.FileInfo) (*models.File, errors.Throwable) {
	of, err := os.Open(path)
	defer of.Close()
	if err != nil {
		return nil, errors.NotFound{}
	}

	if cs.IsForbidden(of.Name()) {
		return nil, errors.AccessDenied{}
	}

	var content []byte
	if loadContent {
		content = make([]byte, fi.Size())
		_, _ = of.Read(content)
	} else {
		content = []byte{}
	}

	mime := cs.GetFileMime(of.Name())
	fm := models.File{
		Name:       cs.TrimName(of.Name()),
		Path:       cs.TrimPath(path),
		Size:       fi.Size(),
		Mime:       mime,
		Hash:       tools.Md5([]byte(path)),
		Permission: fi.Mode().String(),
		Content:    cs.presentableContent(mime, content),
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
	treeName := cs.TrimName(path)
	extra := ""
	if treeName == projectRoot {
		extra = tools.Env("APP_ROOT").(string)
	}
	tree := models.Tree{
		Trees:     []*models.Tree{},
		Files:     []*models.File{},
		Name:      treeName,
		Path:      cs.TrimPath(path),
		Extra:     extra,
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

		f, err := cs.ReadFile(fmt.Sprintf("%s/%s", path, item.Name()), false, item)
		if err != nil {
			continue
		}
		tree.Files = append(tree.Files, f)
	}

	if tree.Trees == nil && tree.Files == nil {
		return nil, nil
	}

	return &tree, nil
}

// GetFileMime returns the mime type of the file based on its extension
func (cs *CodeService) GetFileMime(name string) string {
	list := map[string]string{
		"go":           "application/go",
		"js":           "application/javascript",
		"ts":           "application/typescript",
		"tsx":          "application/typescript-react",
		"jsx":          "application/javascript-react",
		"css":          "text/css",
		"scss":         "text/sass",
		"html":         "text/html",
		"json":         "application/json",
		"jpg":          "image/jpeg",
		"jpeg":         "image/jpeg",
		"png":          "image/png",
		"gif":          "image/gif",
		"svg":          "image/svg+xml",
		"ico":          "image/x-icon",
		"txt":          "text/plain",
		"md":           "text/markdown",
		"sh":           "application/shell",
		"Dockerfile":   "text/dockerfile",
		"dockerfile":   "text/dockerfile",
		"ignore":       "text/ignore",
		"gitignore":    "text/gitignore",
		"dockerignore": "text/dockerignore",
		"yml":          "text/yaml",
	}

	split := strings.Split(name, ".")

	if len(split) >= 1 {
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
		return VirtualRoot
	}

	return fmt.Sprintf("%s/%s", VirtualRoot, p)
}

// TrimName trims out the root path from the file name
func (cs *CodeService) TrimName(name string) string {
	n := strings.Split(cs.TrimPath(name), "/")

	if n[len(n)-1] == "" {
		return projectRoot
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

func (cs *CodeService) presentableContent(mime string, content []byte) string {
	switch mime {
	case "image/jpeg", "image/png", "image/gif", "image/x-icon":
		return fmt.Sprintf("<img src=\"data:image/jpeg;base64,%s\" />", base64.StdEncoding.EncodeToString(content))
	default:
		return string(content)
	}
}
