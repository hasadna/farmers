// The site's server.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

const (
	meteoURL = "http://www.israelmeteo.mobi/Ajax/getStations"
	mainPage = "main.html"
)

var (
	sourceFiles = map[string][]byte{}
	meteoData   []byte

	// File types that should be exposed, along with their MIME types.
	sourceExtensions = map[string]string{
		".html": "text/html",
		".js":   "application/javascript",
		".csv":  "text/plain",
	}
)

func main() {
	log.Print("Loading sources.")
	if err := readSourceFiles("."); err != nil {
		log.Fatal("Failed to load source files: ", err)
	}
	createSourceHandlers()

	log.Print("Getting meteo data.")
	if err := updateMeteoData(); err != nil {
		//log.Fatal("Failed to get meteo data: ", err)
	}
	log.Print(string(meteoData))

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
		if f.IsDir() {
			readSourceFiles(filepath.Join(dir, f.Name()))
			continue
		}
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

		// Main page gets the empty address.
		if f == mainPage {
			f = ""
		}

		http.HandleFunc("/"+f, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", sourceExtensions[filepath.Ext(f)])
			w.Write(data)
		})
	}
}

// updateMeteoData populates meteoData with the latest data from the
// metheorological service.
func updateMeteoData() error {
	res, err := http.Get(meteoURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response status code: %v", res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	meteoData = data
	return nil
}
