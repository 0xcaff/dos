package main

import (
	"log"
	"net/http"
)

// A http.FileSystem for serving single page applications by redirecting all
// unknown paths to index.html.
type SinglePageFileSystem struct {
	backendSystem http.FileSystem
}

func (spa SinglePageFileSystem) Open(name string) (http.File, error) {
	// Try Opening File
	file, err := spa.backendSystem.Open(name)
	if err != nil {
		// Failed to handle opening name, send index.
		log.Printf("[file server] (%s) -> index.html\n", name)
		return spa.backendSystem.Open("index.html")
	} else {
		// Success opening file
		log.Printf("[file server] (%s)\n", name)
		return file, nil
	}
}
