package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// A http.FileSystem for serving single page applications by redirecting all
// unknown paths to index.html.
type SinglePageFileSystem struct {
	backendSystem http.Dir
}

func SPAFileSystem(path string) http.FileSystem {
	return SinglePageFileSystem{
		backendSystem: http.Dir(path),
	}
}

func (spa SinglePageFileSystem) Open(name string) (http.File, error) {
	var localPath string
	basePath := string(spa.backendSystem)
	reqPath := filepath.Join(basePath, name)

	if _, err := os.Stat(reqPath); os.IsNotExist(err) {
		localPath = "index.html"
	} else {
		localPath = name
	}

	fmt.Printf("[file server] (%s) -> (%s)\n", name, localPath)
	return spa.backendSystem.Open(localPath)
}
