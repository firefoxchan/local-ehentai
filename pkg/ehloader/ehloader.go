package ehloader

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type JGallery struct {
	GId int `json:"gid"`
	Token string `json:"token"`
	ArchiverKey string `json:"archiver_key"`
	Title string `json:"title"`
	TitleJpn string `json:"title_jpn"`
	Category string `json:"category"`
	Thumb string `json:"thumb"`
	Uploader string `json:"uploader"`
	Posted json.Number `json:"posted"`
	FileCount json.Number `json:"filecount"`
	FileSize int `json:"filesize"`
	Expunged bool `json:"expunged"`
	Rating json.Number `json:"rating"`
	TorrentCount json.Number `json:"torrentcount"`
	Tags []string `json:"tags"`
}

type Gallery struct {
	GId int
	Token string
	Title string
	TitleJpn string
	Category string
	Thumb string
	Uploader string
	Posted time.Time
	FileCount int
	FileSize int
	Expunged bool
	Rating float32
	TorrentCount int
	Tags map[string][]string
}

type TagK = string
type TagV = string
const (
	TagKAll = "*"
	TagKCategory = "category"
	TagKUploader = "uploader"
	TagKMisc = "misc"
)

var logger = log.New(os.Stderr, "[ehloader]", log.Llongfile | log.LstdFlags)
