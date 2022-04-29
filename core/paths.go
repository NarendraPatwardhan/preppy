package core

import (
	"os"
	"path/filepath"
	"strings"
)

func GetRecursivePaths(root string) ([]string, map[string]bool, error) {
	var paths []string
	folders := make(map[string]bool)
	// Traverse directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Split path by seperator
			pathSegs := strings.Split(path, "/")
			// Get the last segment
			lastSegment := pathSegs[len(pathSegs)-1]
			folders[lastSegment] = true
			return nil
		}
		if filepath.Ext(path) == ".py" {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return paths, folders, nil
}
