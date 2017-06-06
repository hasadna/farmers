// The site's server.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const (
	meteoURL             = "http://www.israelmeteo.mobi/Ajax/getStations"
	mainPage             = "main.html"
	meteoUpdateFrequency = 10 * time.Minute
)

var (
	sourceFiles = map[string][]byte{}
	meteoData   []byte
	debug       = false

	// File types that should be exposed, along with their MIME types.
	sourceExtensions = map[string]string{
		".html": "text/html",
		".js":   "application/javascript",
		".csv":  "text/plain",
	}
)

func main() {
	parseFlags()

	// Load sources.
	if debug {
		log.Print("Running in debug mode; serving files dynamically.")
		http.HandleFunc("/"+mainPage, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			page, _ := ioutil.ReadFile(mainPage)
			w.Write(page)
		})
		http.Handle("/", http.FileServer(http.Dir(".")))
	} else {
		log.Print("Loading sources.")
		if err := readSourceFiles("."); err != nil {
			log.Fatal("Failed to load source files: ", err)
		}
		createSourceHandlers()
	}

	// Get meteo data.
	log.Print("Getting meteo data.")
	if err := updateMeteoData(); err != nil {
		log.Fatal("Failed to get meteo data: ", err)
	}
	http.HandleFunc("/meteodata.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(meteoData)
	})
	go updateMeteoDataPeriodically()

	// Start server!
	log.Print("Listening.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseFlags() {
	flag.BoolVar(&debug, "debug", false,
		"Development mode - does not load sources in advance, but serves them dynamically.")
	flag.Parse()
}

// readSourceFiles reads the source files and places their data in sourceFiles.
func readSourceFiles(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory content: %s", err)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
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
		sourceFiles[filepath.Join(dir, f.Name())] = data
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

// updateMeteoDataPeriodically updates the meteoData variable.
func updateMeteoDataPeriodically() {
	for {
		time.Sleep(meteoUpdateFrequency)
		err := updateMeteoData()
		if err != nil {
			log.Print("Failed to update meteo data: ", err)
		} else {
			log.Print("Updated meteo data.")
		}
	}
}
