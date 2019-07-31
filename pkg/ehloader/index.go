package ehloader

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

var tags map[TagK]map[TagV][]int
var galleries map[int]*Gallery
var indexMu sync.RWMutex
var threadNum = 8


func ScanJson(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	inColumn:=false
	level,start,end:=0,0,0
	for l,s := range data{
		if s == '"'{
			inColumn = !inColumn
		}
		if s == '{' &&!inColumn{
			if start == 0{
				start = l
			}
			level++
		}
		if s == '}' &&!inColumn{
			level--
			if level <= 0{
				end = l
				break
			}
		}
	}
	if end > start{
		return end+1, data[start:end+1], nil
		// We have a full newline-terminated line.
	}//logger.Println(start,end,level,len(data))

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), nil, nil
	}
	// Request more data.
	return 0, nil, nil
}

func feedJson(feedCh chan string,jsonCh chan JGallery,barrier *sync.WaitGroup){
	j := JGallery{}
	for b := range feedCh{
		e := json.Unmarshal([]byte(b),&j)
		if e != nil{
			logger.Println(e)
			continue
		}
		jsonCh <- j
	}
	barrier.Done()
}

func tagJson(jsonCh chan JGallery,barrier *sync.WaitGroup){
	for j := range jsonCh{
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
		//delete(jGalleries, i)
	}
	barrier.Done()
}

func IndexJson(path string) error {
	indexMu.Lock()
	defer indexMu.Unlock()
	logger.Printf("Start Parsing json.\n")
	feedCh := make(chan string,2*threadNum)
	jsonCh := make(chan JGallery,2*threadNum)
	//jGalleries := make(map[int64]JGallery, 850000)
	feedBarrier := sync.WaitGroup{}
	feedBarrier.Add(threadNum)
	tagBarrier := sync.WaitGroup{}
	tagBarrier.Add(1)
	galleries = map[int]*Gallery{}
	tags = map[TagK]map[TagV][]int{}
	for i:=0;i<threadNum;i++{
		go feedJson(feedCh,jsonCh,&feedBarrier)
	}
	go tagJson(jsonCh,&tagBarrier)

	f, e := os.Open(path)
	f.Seek(1,0)
	//skip first {
	if e != nil {
		return e
	}

	b := bufio.NewScanner(f)
	b.Split(ScanJson)

	count := 0
	for b.Scan(){
		feedCh <- b.Text()
		count++
		//logger.Println(j)
		if count%10000 == 0{
			//logger.Println(j)
			logger.Println(count)
		}
	}
	close(feedCh)
	f.Close()
	//logger.Println(j.GId)
	logger.Println(count)


	logger.Printf("End Parsing json.\n")
	logger.Printf("Start Loading gallaries.\n")
	feedBarrier.Wait()
	close(jsonCh)
	tagBarrier.Wait()


	//for _, tagVs := range tags {
	//	for value := range tagVs {
	//		sort.Ints(tagVs[value])
	//	}
	//}
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
