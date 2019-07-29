package ehloader

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var tags map[TagK]map[TagV][]int
var galleries map[int]*Gallery
var indexMu sync.RWMutex

func IndexJson(path string) error {
	indexMu.Lock()
	defer indexMu.Unlock()
	logger.Printf("Start Parsing json.\n")
	jGalleries := make(map[int64]JGallery, 850000)
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
	tags = map[TagK]map[TagV][]int{}
	logger.Printf("Start Loading gallaries.\n")
	counter := 0
	for i, j := range jGalleries {
		if counter % 5000 == 0 {
			logger.Printf("Loading %d gallaries...\n", counter)
		}
		counter++
		gallery := &Gallery{
			GId:          j.GId,
			Token:        strings.TrimSpace(j.Token),
			Title:        strings.TrimSpace(j.Title),
			TitleJpn:     strings.TrimSpace(j.TitleJpn),
			Category:     strings.TrimSpace(j.Category),
			Thumb:        j.Thumb,
			Uploader:     strings.TrimSpace(j.Uploader),
			Posted:       time.Time{},
			FileCount:    0,
			FileSize:     j.FileSize,
			Expunged:     j.Expunged,
			Rating:       0,
			TorrentCount: 0,
			Tags:         map[TagK][]TagV{},
		}
		{
			posted, e := j.Posted.Int64()
			if e != nil {
				logger.Printf("Load gallary %d .posted error: %s\n", j.GId, e)
			} else {
				gallery.Posted = time.Unix(posted, 0)
			}
		}
		{
			fc, e := j.FileCount.Int64()
			if e != nil {
				logger.Printf("Load gallary %d .filecount error: %s\n", j.GId, e)
			} else {
				gallery.FileCount = int(fc)
			}
		}
		{
			rt, e := j.Rating.Float64()
			if e != nil {
				logger.Printf("Load gallary %d .rating error: %s\n", j.GId, e)
			} else {
				gallery.Rating = float32(rt)
			}
		}
		{
			tc, e := j.TorrentCount.Int64()
			if e != nil {
				logger.Printf("Load gallary %d .torrent_count error: %s\n", j.GId, e)
			} else {
				gallery.TorrentCount = int(tc)
			}
		}
		for _, pair := range j.Tags {
			key, value := BuildKV(pair, TagKMisc)
			appendTagKVG(key, value, j.GId)
			if _, ok := gallery.Tags[key]; !ok {
				gallery.Tags[key] = make([]string, 0)
			}
			gallery.Tags[key] = append(gallery.Tags[key], value)
		}
		appendTagKVG(TagKCategory, gallery.Category, j.GId)
		appendTagKVG(TagKUploader, gallery.Uploader, j.GId)
		galleries[j.GId] = gallery
		delete(jGalleries, i)
	}
	for _, tagVs := range tags {
		for value := range tagVs {
			sort.Ints(tagVs[value])
		}
	}
	logger.Printf("Loading %d gallaries...\n", counter)
	logger.Printf("End Loading gallaries.\n")
	return nil
}

func BuildKV (pair string, defaultTagK string) (string, string) {
	pairs := strings.SplitN(pair, ":", 2)
	var key, value string
	switch len(pairs) {
	case 1:
		key = defaultTagK
		value = strings.TrimSpace(pairs[0])
	case 2:
		key = strings.TrimSpace(pairs[0])
		value = strings.TrimSpace(pairs[1])
	}
	return key, value
}

func appendTagKVG (key, value string, gid int) {
	key = strings.ToLower(key)
	value = strings.ToLower(value)
	if _, ok := tags[key]; !ok {
		tags[key] = map[TagV][]int{}
	}
	if _, ok := tags[key][value]; !ok {
		tags[key][value] = make([]int, 0)
	}
	tags[key][value] = append(tags[key][value], gid)
}
