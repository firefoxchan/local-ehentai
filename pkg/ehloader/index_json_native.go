package ehloader

import (
	"encoding/json"
	"os"
)

func indexJsonNative(path string) error {
	logger.Printf("Start Parsing json.\n")
	jGalleries := make(map[int]JGallery, 850000)
	{
		f, e := os.OpenFile(path, os.O_RDONLY, 0)
		if e != nil {
			return e
		}
		dec := json.NewDecoder(f)
		if e := dec.Decode(&jGalleries); e != nil {
			_ = f.Close()
			return e
		}
		_ = f.Close()
	}
	logger.Printf("End Parsing json.\n")
	galleries = make(map[int]*Gallery, 850000)
	logger.Printf("Start Loading gallaries.\n")
	counter := 0
	for i, j := range jGalleries {
		if counter%10000 == 0 {
			logger.Printf("Loading %d gallaries...\n", counter)
		}
		counter++
		handleJGallery(j)
		delete(jGalleries, i)
	}
	logger.Printf("Loading %d gallaries...\n", counter)
	logger.Printf("End Loading gallaries.\n")
	return nil
}
