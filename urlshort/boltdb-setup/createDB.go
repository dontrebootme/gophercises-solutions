package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type URLShort struct {
	ID   int
	Path string
	URL  string
}

func main() {
	db, err := bolt.Open("urlshort.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var urlshorts []URLShort
	urlshorts = []URLShort{
		{
			Path: "/patrick1",
			URL:  "https://github.com/dontrebootme",
		},
		{
			Path: "/disney",
			URL:  "https://disney.com",
		},
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("urlshort"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	for k, _ := range urlshorts {
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("urlshort"))
			id, _ := b.NextSequence()
			urlshorts[k].ID = int(id)
			buf, err := json.Marshal(urlshorts[k])
			if err != nil {
				return err
			}
			b.Put(itob(urlshorts[k].ID), buf)
			return nil
		})
	}

}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
