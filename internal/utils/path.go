package utils

import (
	"path/filepath"
	"strings"
)

// NormalizePath ensures the path has a trailing slash if it's a directory
func NormalizePath(filePath string, isDirectory bool) string {
	if !filepath.IsAbs(filePath) && !strings.HasPrefix(filePath, ".") {
		filePath = "./" + filePath
	}

	if isDirectory && !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}

	return filepath.Clean(filePath)
}
