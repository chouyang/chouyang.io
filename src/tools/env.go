package tools

import (
	"os"
	"strings"
)

// configs is a map of environment variables to their values.
var configs map[string]interface{}

func init() {
	if configs != nil {
		return
	}

	if fi, err := os.Stat(".env"); err == nil {
		loadEnv(".env", fi.Size())
	}
}

// Env returns the value of the environment variable.
func Env(path string) interface{} {
	if val, ok := configs[path]; ok {
		return val
	}

	return nil
}

// loadEnv loads the environment variables into the configs map.
func loadEnv(path string, size int64) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return
	}

	buf := make([]byte, size)
	_, _ = file.Read(buf)
	lines := strings.Split(string(buf), "\n")
	configs = make(map[string]interface{}, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			kv := strings.Split(line, "=")
			configs[kv[0]] = kv[1]
		}
	}

	return
}
