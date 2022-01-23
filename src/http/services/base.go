package services

import (
	"chouyang.io/src/tools"
	"os"
	"path"
)

// RealRoot points to the parent path of the current working directory
var RealRoot string

// VirtualRoot represents the virtual root path of the project
var VirtualRoot string

// ignoreList is a list of files and directories to ignore
var ignoreList []string

// init initializes the package variables
func init() {
	RealRoot, _ = os.Getwd()
	RealRoot = path.Dir(RealRoot)

	VirtualRoot = tools.Env("APP_ROOT").(string)

	ignoreList = []string{".env", ".DS_Store", ".git/", ".idea/", ".vscode/", "node_modules/"}
}
