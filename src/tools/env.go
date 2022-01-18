package tools

import (
	"os"
	"strings"
)

var env map[string]interface{}

func init() {
	if env != nil {
		return
	}

	if fi, err := os.Stat(".env"); err == nil {
		loadEnv(".env", fi.Size())
	}
}

func Env(path string) interface{} {
	if val, ok := env[path]; ok {
		return val
	}

	return nil
}

func loadEnv(path string, size int64) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return
	}

	buf := make([]byte, size)
	_, _ = file.Read(buf)
	lines := strings.Split(string(buf), "\n")
	env = make(map[string]interface{}, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			kv := strings.Split(line, "=")
			env[kv[0]] = kv[1]
		}
	}

	return
}
