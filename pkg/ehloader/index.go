package ehloader

import (
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var tags map[TagK]map[TagV][]int
var galleries map[int]*Gallery
var gExistsInUrls map[int]struct{}
var indexMu sync.RWMutex

func Index(jsonPath string, urlPath string) error {
	indexMu.Lock()
	defer indexMu.Unlock()
	if urlPath != "" {
		matchedGs, e := indexURLList(urlPath)
		if e != nil {
			return e
		}
		gExistsInUrls = matchedGs
	}
	return indexJsonFast(jsonPath)
}

func handleJGallery(j JGallery) {
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
	// tags
	for _, pair := range j.Tags {
		key, value := BuildKV(pair, TagKMisc)
		appendTagKVG(key, value, j.GId)
		if _, ok := gallery.Tags[key]; !ok {
			gallery.Tags[key] = make([]string, 0)
		}
		gallery.Tags[key] = append(gallery.Tags[key], value)
	}
	// category / uploader
	appendTagKVG(TagKCategory, gallery.Category, j.GId)
	appendTagKVG(TagKUploader, gallery.Uploader, j.GId)
	// expunged
	switch gallery.Expunged {
	case true:
		appendTagKVG(TagKExpunged, TagVExpungedTrue, j.GId)
	case false:
		appendTagKVG(TagKExpunged, TagVExpungedFalse, j.GId)
	}
	// min rating
	for i := int64(0); i <= int64(math.Round(float64(gallery.Rating))); i++ {
		appendTagKVG(TagKMinRating, strconv.FormatInt(i, 10), j.GId)
	}
	// exists
	switch {
	case existsIn(j.GId, gExistsInUrls):
		gallery.Exists = true
		appendTagKVG(TagKExists, TagVExistsTrue, j.GId)
		gallery.ExistsIn = TagVExistsInURL
		appendTagKVG(TagKExistsIn, TagVExistsInURL, j.GId)
	default:
		gallery.Exists = false
		appendTagKVG(TagKExists, TagVExistsFalse, j.GId)
	}
	galleries[j.GId] = gallery
}

func existsIn(gid int, set map[int]struct{}) bool {
	_, ok := set[gid]
	return ok
}

func BuildKV(pair string, defaultTagK string) (string, string) {
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

func appendTagKVG(key, value string, gid int) {
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
