package urlshort

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	yamlV2 "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to aim any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYAML)
	return MapHandler(pathMap, fallback), nil
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

type urlShortRow struct {
	Path string
	URL  string
}

func BoltDBHandler(boltDBFile string, fallback http.Handler) (http.HandlerFunc, error) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(boltDBFile, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pathMap := make(map[string]string)
	db.View(func(tx *bolt.Tx) error {

		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("urlshort"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var row urlShortRow
			err := json.Unmarshal(v, &row)
			if err != nil {
				panic(err)
			}
			//			fmt.Printf("key=%s, value=%s\n", k, v)
			pathMap[row.Path] = row.URL
		}
		return nil
	})

	//	pathMap := buildMap(parsedJSON)
	//	return MapHandler(pathMap, fallback), nil
	return MapHandler(pathMap, fallback), nil

}

func parseYAML(yaml []byte) (dst []map[string]string, err error) {
	err = yamlV2.Unmarshal(yaml, &dst)
	return dst, err
}

func parseJSON(j []byte) (dst []map[string]string, err error) {
	err = json.Unmarshal(j, &dst)
	return dst, err
}

func buildMap(parsed []map[string]string) map[string]string {
	mergedMap := make(map[string]string)
	for _, entry := range parsed {
		key := entry["path"]
		mergedMap[key] = entry["url"]
	}
	return mergedMap
}
