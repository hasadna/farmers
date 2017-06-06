// The site's server.
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

const (
	meteoURL = "www.israelmeteo.mobi/Ajax/getStations"
)

var (
	sourceFiles = map[string][]byte{}

	// File types that should be exposed, along with their MIME types.
	sourceExtensions = map[string]string{
		".html": "text/html",
		".js":   "application/javascript",
		".csv":  "text/plain",
	}
)

func main() {
	log.Print("Reading sources.")
	if err := readSourceFiles("."); err != nil {
		log.Fatal("Failed to load source files: ", err)
	}
	createSourceHandlers()

	log.Print("Listening.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// readSourceFiles reads the source files and places their data in sourceFiles.
func readSourceFiles(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if sourceExtensions[filepath.Ext(f.Name())] == "" {
			continue
		}
		data, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return err
		}
		sourceFiles[f.Name()] = data
	}
	return nil
}

// createSourceHandlers creates an HTTP handler for each source file.
func createSourceHandlers() {
	for f, data := range sourceFiles {
		// Create local context for anonymous function.
		f := f
		data := data
		log.Print(f)
		http.HandleFunc("/"+f, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", sourceExtensions[filepath.Ext(f)])
			w.Write(data)
		})
	}
}
