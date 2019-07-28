package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/firefoxchan/local-ehentai/pkg/ehloader"
	"os"
	"strconv"
	"strings"
)

func main ()  {
	var (
		jsonPath string
		format string
	)
	flag.StringVar(&jsonPath, "j", "", "path to eh api json file")
	flag.StringVar(&format, "f", "dense", "output format. dense, json")
	flag.Parse()
	if jsonPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if e := ehloader.IndexJson(jsonPath); e != nil {
		panic(e)
	}
	// search
	scanner := bufio.NewScanner(os.Stdin)
	printHint := func () {
		fmt.Printf("Search format: offset, limit, tag1:value1, tag2:value2, ...\n")
		fmt.Printf("Example: 0, 10, category:doujinshi, parody:the idolmaster\n")
		fmt.Printf("Quit format: quit\n")
		fmt.Printf("Please Input: ")
	}
	printHint()
	for scanner.Scan() {
		in := scanner.Text()
		switch in {
		case "exit", "quit":
			fmt.Printf("Bye\n")
			return
		}
		params := strings.SplitN(in, ",", 3)
		if len(params) != 3 {
			printHint()
			continue
		}
		offset, e := strconv.ParseInt(strings.TrimSpace(params[0]), 10, 64)
		if e != nil {
			fmt.Println("Unable to parse offset:", e)
			printHint()
			continue
		}
		limit, e := strconv.ParseInt(strings.TrimSpace(params[1]), 10, 64)
		if e != nil {
			fmt.Println("Unable to parse limit:", e)
			printHint()
			continue
		}
		search := strings.TrimSpace(params[2])
		kvs := strings.Split(search, ",")
		searchTags := map[ehloader.TagK]map[ehloader.TagV]struct{}{}
		for _, kv := range kvs {
			key, value := ehloader.BuildKV(kv, ehloader.TagKAll)
			if _, ok := searchTags[key]; !ok {
				searchTags[key] = map[ehloader.TagV]struct{}{}
			}
			searchTags[key][value] = struct{}{}
		}
		ret := ehloader.Search(searchTags, int(offset), int(limit))
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
		printHint()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
