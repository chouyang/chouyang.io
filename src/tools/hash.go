package tools

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// Md5 returns the md5 hash of the giving string of bytes
func Md5(content []byte) string {
	hash := md5.New()
	hash.Write(content)

	return strings.ToUpper(fmt.Sprintf("%x", hash.Sum(nil)))
}
