package main

import (
	log "github.com/sirupsen/logrus"
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
		log.Infof("[file server] (%s) -> index.html", name)
		return spa.backendSystem.Open("index.html")
	} else {
		// Success opening file
		log.Infof("[file server] (%s)", name)
		return file, nil
	}
}
