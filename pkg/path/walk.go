package path

import (
	"os"
	"path/filepath"
	"strings"
)

func Walk(workPath string, f func(filePath string) error) {
	err := filepath.Walk(workPath, func(filePath string, fileInfo os.FileInfo, _ error) error {
		if !fileInfo.IsDir() && strings.HasSuffix(filePath, ".go") && !strings.HasSuffix(filePath, "_test.go") {
			return f(filePath)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
