package ehloader

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var tags map[TagK]map[TagV][]int
var galleries map[int]*Gallery
var gExistsInUrls map[int]struct{}
var rIdxTitle map[string]map[string][]int
var rIdxTitleJpn map[string]map[string][]int
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
		ArchiverKey:  j.ArchiverKey,
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
		posted, e := strconv.ParseInt(j.Posted, 10, 64)
		if e != nil {
			logger.Printf("Load gallary %d .posted error: %s\n", j.GId, e)
		} else {
			gallery.Posted = time.Unix(posted, 0)
		}
	}
	{
		fc, e := strconv.ParseInt(j.FileCount, 10, 64)
		if e != nil {
			logger.Printf("Load gallary %d .filecount error: %s\n", j.GId, e)
		} else {
			gallery.FileCount = int(fc)
		}
	}
	{
		rt, e := strconv.ParseFloat(j.Rating, 64)
		if e != nil {
			logger.Printf("Load gallary %d .rating error: %s\n", j.GId, e)
		} else {
			gallery.Rating = float32(rt)
		}
	}
	{
		tc, e := strconv.ParseInt(j.TorrentCount, 10, 64)
		if e != nil {
			logger.Printf("Load gallary %d .torrent_count error: %s\n", j.GId, e)
		} else {
			gallery.TorrentCount = int(tc)
		}
	}
	gallery.TitleExt = parseTitle(gallery.Title)
	gallery.TitleJpnExt = parseTitle(gallery.TitleJpn)
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

var titleParseRes = []*regexp.Regexp{
	regexp.MustCompile(`(?P<convention>\([^)]+\))?[ ]*(?P<groupArtist>\[[^]]+])?(?P<title>[^([]*)(?P<parody>\([^)]+\))?[ ]*(?P<translation>\[[^]]+])?`),
}

const (
	rIdxTitleConvention  = "Convention"
	rIdxTitleGroup       = "Group"
	rIdxTitleArtist      = "Artist"
	rIdxTitleTitle       = "Title"
	rIdxTitleParody      = "Parody"
	rIdxTitleTranslation = "Translation"
)
var rIdxTitleAll = []string{
	rIdxTitleConvention,
	rIdxTitleGroup,
	rIdxTitleArtist,
	rIdxTitleTitle,
	rIdxTitleParody,
	rIdxTitleTranslation,
}

var titleSPCharReplacer = strings.NewReplacer(
	`（`, `(`, `）`, `)`, `【`, `[`, `】`, `]`,
	`（`, `(`, `）`, `)`, `［`, `[`, `］`, `]`,
)

func parseTitle (title string) map[string]string {
	title = titleSPCharReplacer.Replace(title)
	t := map[string]string{}
	for _, re := range titleParseRes {
		groupNames := re.SubexpNames()
		for _, match := range re.FindAllStringSubmatch(title, -1) {
			for groupIdx, matched := range match {
				name := groupNames[groupIdx]
				if name == "" {
					name = "*"
				}
				if _, ok := t[name]; ok {
					continue
				}
				matched := strings.ToLower(strings.TrimSpace(matched))
				if matched == "" {
					continue
				}
				switch name {
				case "convention":
					t[rIdxTitleConvention] = strings.Trim(matched, "() ")
				case "title":
					t[rIdxTitleTitle] = matched
				case "parody":
					t[rIdxTitleParody] = strings.Trim(matched, "() ")
				case "translation":
					t[rIdxTitleTranslation] = strings.Trim(matched, "[] ")
				case "groupArtist":
					matched = strings.Trim(matched, "[] ")
					if strings.Contains(matched, "(") && strings.HasSuffix(matched, ")") {
						// group (artist)
						matched = strings.TrimSuffix(matched, ")")
						pairs := strings.SplitN(matched, "(", 2)
						t[rIdxTitleGroup] = strings.TrimSpace(pairs[0])
						t[rIdxTitleArtist] = strings.TrimSpace(pairs[1])
					} else {
						t[rIdxTitleArtist] = matched
					}
				}
			}
		}
	}
	for _, rIdx := range rIdxTitleAll {
		if _, ok := t[rIdx]; !ok {
			t[rIdx] = ""
		}
	}
	return t
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
