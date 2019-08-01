package ehloader

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func IndexJsonNative(path string) error {
	indexMu.Lock()
	defer indexMu.Unlock()
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
		if counter % 10000 == 0 {
			logger.Printf("Loading %d gallaries...\n", counter)
		}
		counter++
		handleJGallery(j)
		delete(jGalleries, i)
	}
	logger.Printf("Loading %d gallaries...\n", counter)
	tagDumps := make([]string, 0)
	const tagDumpMinLength = 100
	for tagK, tagVs := range tags {
		tagKDumps := make([]string, 0)
		for value := range tagVs {
			sort.Ints(tagVs[value])
			if len(tagVs[value]) > tagDumpMinLength {
				tagKDumps = append(tagKDumps, fmt.Sprintf("%s:%d", value, len(tagVs[value])))
			}
		}
		switch tagK {
		case "male", "female", "misc", "language", TagKCategory, TagKMinRating, TagKExpunged:
			// pass
		default:
			continue
		}
		tagDumps = append(tagDumps, fmt.Sprintf("  %s:", tagK))
		const oneLineLimit = 10
		for oneLineLimit < len(tagKDumps) {
			tagKDumps, tagDumps = tagKDumps[oneLineLimit:], append(tagDumps, fmt.Sprintf("    %s", strings.Join(tagKDumps[0:oneLineLimit:oneLineLimit], ", ")))
		}
		tagDumps = append(tagDumps, fmt.Sprintf("    %s", strings.Join(tagKDumps, ", ")))
	}
	logger.Printf("Tag stats (> %d):\n%s", tagDumpMinLength, strings.Join(tagDumps, "\n"))
	logger.Printf("End Loading gallaries.\n")
	return nil
}
