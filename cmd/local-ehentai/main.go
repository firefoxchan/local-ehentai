package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/firefoxchan/local-ehentai/pkg/ehloader"
	"os"
	"strings"
)

func main ()  {
	var (
		jsonPath string
		search string
		offset int
		limit int
		format string
	)
	flag.StringVar(&jsonPath, "j", "", "path to eh api json file")
	flag.StringVar(&search, "s", "", "search string, eg: \"category:doujinshi, parody:the idolmaster\"")
	flag.IntVar(&offset, "o", 0, "search offset")
	flag.IntVar(&limit, "l", 10, "search limit")
	flag.StringVar(&format, "f", "dense", "output format. dense, json")
	flag.Parse()
	if jsonPath == "" || search == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if e := ehloader.IndexJson(jsonPath); e != nil {
		panic(e)
	}
	// search
	kvs := strings.Split(search, ",")
	searchTags := map[ehloader.TagK]map[ehloader.TagV]struct{}{}
	for _, kv := range kvs {
		key, value := ehloader.BuildKV(kv, ehloader.TagKAll)
		if _, ok := searchTags[key]; !ok {
			searchTags[key] = map[ehloader.TagV]struct{}{}
		}
		searchTags[key][value] = struct{}{}
	}
	ret := ehloader.Search(searchTags, offset, limit)
	fmt.Printf("Search: %d, %d, %v -> %v\n", offset, limit, search, searchTags)
	if format == "dense" {
		for _, g := range ret {
			b, e := json.Marshal(g)
			if e != nil {
				fmt.Printf("  %+v\n", e)
			} else {
				fmt.Printf("  %s\n", string(b))
			}
		}
	} else {
		b, e := json.MarshalIndent(ret, "  ", "  ")
		if e != nil {
			fmt.Printf("  %+v\n", e)
		} else {
			fmt.Printf("  %s\n", string(b))
		}
	}
}
