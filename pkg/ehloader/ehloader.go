package ehloader

import (
	"log"
	"os"
	"time"
)

type JGallery struct {
	GId          int      `json:"gid"`
	Token        string   `json:"token"`
	ArchiverKey  string   `json:"archiver_key"`
	Title        string   `json:"title"`
	TitleJpn     string   `json:"title_jpn"`
	Category     string   `json:"category"`
	Thumb        string   `json:"thumb"`
	Uploader     string   `json:"uploader"`
	Posted       string   `json:"posted"`
	FileCount    string   `json:"filecount"`
	FileSize     int64    `json:"filesize"`
	Expunged     bool     `json:"expunged"`
	Rating       string   `json:"rating"`
	TorrentCount string   `json:"torrentcount"`
	Tags         []string `json:"tags"`
}

type Gallery struct {
	GId          int
	Token        string
	ArchiverKey  string
	Title        string
	TitleExt     map[string]string
	TitleJpn     string
	TitleJpnExt  map[string]string
	Category     string
	Thumb        string
	Uploader     string
	Posted       time.Time
	FileCount    int
	FileSize     int64
	Expunged     bool
	Rating       float32
	TorrentCount int
	Tags         map[string][]string
	Exists       bool
	ExistsIn     string
	ExistFiles   []string
}

type TagK = string
type TagV = string

const (
	TagKAll       TagK = "*"
	TagKCategory  TagK = "category"
	TagKUploader  TagK = "uploader"
	TagKMisc      TagK = "misc"
	TagKExpunged  TagK = "expunged"
	TagKMinRating TagK = "min rating"
	TagKExists    TagK = "exists"
	TagKExistsIn  TagK = "exists in"
	TagKGId       TagK = "gid"
)

// rIdx
const (
	// title
	TagKRIdxTitlePrefix      TagK = "title"
	TagKRIdxTitleConvention  TagK = TagKRIdxTitlePrefix + " " + rIdxTitleConvention
	TagKRIdxTitleGroup       TagK = TagKRIdxTitlePrefix + " " + rIdxTitleGroup
	TagKRIdxTitleArtist      TagK = TagKRIdxTitlePrefix + " " + rIdxTitleArtist
	TagKRIdxTitleTitle       TagK = TagKRIdxTitlePrefix + " " + rIdxTitleTitle
	TagKRIdxTitleParody      TagK = TagKRIdxTitlePrefix + " " + rIdxTitleParody
	TagKRIdxTitleTranslation TagK = TagKRIdxTitlePrefix + " " + rIdxTitleTranslation
	// title jpn
	TagKRIdxTitleJpnPrefix      TagK = "title_jpn"
	TagKRIdxTitleJpnConvention  TagK = TagKRIdxTitleJpnPrefix + " " + rIdxTitleConvention
	TagKRIdxTitleJpnGroup       TagK = TagKRIdxTitleJpnPrefix + " " + rIdxTitleGroup
	TagKRIdxTitleJpnArtist      TagK = TagKRIdxTitleJpnPrefix + " " + rIdxTitleArtist
	TagKRIdxTitleJpnTitle       TagK = TagKRIdxTitleJpnPrefix + " " + rIdxTitleTitle
	TagKRIdxTitleJpnParody      TagK = TagKRIdxTitleJpnPrefix + " " + rIdxTitleParody
	TagKRIdxTitleJpnTranslation TagK = TagKRIdxTitleJpnPrefix + " " + rIdxTitleTranslation
)

const (
	TagVExpungedFalse TagV = "0"
	TagVExpungedTrue  TagV = "1"
	TagVExistsFalse   TagV = "0"
	TagVExistsTrue    TagV = "1"
	TagVExistsInURL   TagV = "url"
	TagVExistsInFile  TagV = "file"
)

var logger = log.New(os.Stderr, "[ehloader]", log.Lshortfile|log.LstdFlags)
