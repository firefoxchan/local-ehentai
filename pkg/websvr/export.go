package websvr

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/firefoxchan/local-ehentai/pkg/ehloader"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func exportCSV(pq parsedQuery, writer http.ResponseWriter, request *http.Request) {
	gs, _ := ehloader.SearchQ(0, -1, pq.Q)
	writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="export.%s.csv"`, time.Now().Format("20060102150405")))
	// compatible with excel
	if _, e := writer.Write([]byte{'\xef', '\xbb', '\xbf'}); e != nil {
		logger.Printf("Export csv write error: %s\n", e)
		return
	}
	wt := csv.NewWriter(writer)
	columns := []string{
		"gid", "token", "link",
		"category", "uploader", "filecount", "filesize",
		"expunged", "posted",
		"rating",
		"title", "title_jpn", "thumb",
		"archiver_key",
		"tags",
	}
	if e := wt.Write(columns); e != nil {
		logger.Printf("Export csv write error: %s\n", e)
		return
	}
	counter := 0
	for _, g := range gs {
		data := []string{
			strconv.FormatInt(int64(g.GId), 10), g.Token, fmt.Sprintf("https://e-hentai.org/g/%d/%s", g.GId, g.Token),
			g.Category, g.Uploader, strconv.FormatInt(int64(g.FileCount), 10), strconv.FormatInt(int64(g.FileSize), 10),
			strconv.FormatBool(g.Expunged), g.Posted.Format("2006-01-02 15:04:05"),
			strconv.FormatFloat(float64(g.Rating), 'f', 2, 64),
			g.Title, g.TitleJpn, g.Thumb,
			g.ArchiverKey,
		}
		tags := make([]string, 0)
		for key, values := range g.Tags {
			tags = append(tags, fmt.Sprintf("%s:%s", key, strings.Join(values, ",")))
		}
		data = append(data, strings.Join(tags, ";"))
		if e := wt.Write(data); e != nil {
			logger.Printf("Export csv write error: %s\n", e)
			return
		}
		if counter%10000 == 0 {
			logger.Printf("Exported %d gallaries.\n", counter)
		}
		counter++
		if counter%100 == 0 {
			wt.Flush()
			if e := wt.Error(); e != nil {
				logger.Printf("Export csv write error: %s\n", e)
				return
			}
		}
	}
	wt.Flush()
	if e := wt.Error(); e != nil {
		logger.Printf("Export csv write error: %s\n", e)
		return
	}
	logger.Printf("Exported %d gallaries.\n", counter)
}

func exportJSON(pq parsedQuery, writer http.ResponseWriter, request *http.Request) {
	gs, _ := ehloader.SearchQ(0, -1, pq.Q)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="export.%s.json"`, time.Now().Format("20060102150405")))
	if _, e := writer.Write([]byte{'{'}); e != nil {
		logger.Printf("Export json write error: %s\n", e)
		return
	}
	counter := 0
	for _, g := range gs {
		jG := ehloader.JGallery{
			GId:          g.GId,
			Token:        g.Token,
			ArchiverKey:  g.ArchiverKey,
			Title:        g.Title,
			TitleJpn:     g.TitleJpn,
			Category:     g.Category,
			Thumb:        g.Thumb,
			Uploader:     g.Uploader,
			Posted:       strconv.FormatInt(g.Posted.Unix(), 10),
			FileCount:    strconv.FormatInt(int64(g.FileCount), 10),
			FileSize:     g.FileSize,
			Expunged:     g.Expunged,
			Rating:       strconv.FormatFloat(float64(g.Rating), 'f', 2, 32),
			TorrentCount: strconv.FormatInt(int64(g.TorrentCount), 10),
			Tags:         []string{},
		}
		for key, values := range g.Tags {
			for _, value := range values {
				jG.Tags = append(jG.Tags, fmt.Sprintf("%s:%s", key, value))
			}
		}
		if counter != 0 {
			if _, e := writer.Write([]byte{','}); e != nil {
				logger.Printf("Export json write error: %s\n", e)
				return
			}
		}
		b, e := json.Marshal(jG)
		if e != nil {
			logger.Printf("Export json marshal error: %s\n", e)
			return
		}
		if _, e := fmt.Fprintf(writer, `"%d":%s`, g.GId, b); e != nil {
			logger.Printf("Export json write error: %s\n", e)
			return
		}
		if counter%10000 == 0 {
			logger.Printf("Exported %d gallaries.\n", counter)
		}
		counter++
	}
	if _, e := writer.Write([]byte{'}'}); e != nil {
		logger.Printf("Export json write error: %s\n", e)
		return
	}
	logger.Printf("Exported %d gallaries.\n", counter)
}
