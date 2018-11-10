package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dontrebootme/gophercises/urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/dontrebootme":   "https://dontreboot.me/",
		"/spof.io":        "https://spof.io/",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	var yamlFile = flag.String("yaml", "urlshort.yaml", "yaml file for url short paths")
	var jsonFile = flag.String("json", "urlshort.json", "json file for url short paths")
	var boltDBFile = flag.String("db", "urlshort.db", "bolt db file for url short paths")
	flag.Parse()

	yFile, err := ioutil.ReadFile(*yamlFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to open the yaml file: %s\n", *yamlFile))
	}
	jFile, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to open the json file: %s\n", *jsonFile))
	}

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(yFile, mapHandler)
	if err != nil {
		panic(err)
	}
	// Build the JSONHandler using the yamlHandler as the
	// fallback
	jsonHandler, err := urlshort.JSONHandler(jFile, yamlHandler)
	if err != nil {
		panic(err)
	}

	boltDBHandler, err := urlshort.BoltDBHandler(*boltDBFile, jsonHandler)

	fmt.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", boltDBHandler); err != nil {
		panic(err)
	}

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
