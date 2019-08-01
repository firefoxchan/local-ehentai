package ehloader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// Buggy
func IndexJsonFast(path string) error {
	indexMu.Lock()
	defer indexMu.Unlock()

	var feederNum = int(math.Max(float64(runtime.NumCPU()-2), 1))
	rawJsonCh := make(chan string, 2*feederNum)
	jGalleriesCh := make(chan JGallery, 2*feederNum)

	logger.Printf("Start Parsing json (%d gr).\n", feederNum)
	feedBarrier := sync.WaitGroup{}
	feedBarrier.Add(feederNum)
	for i := 0; i < feederNum; i++ {
		go feedJson(rawJsonCh, jGalleriesCh, &feedBarrier)
	}

	logger.Printf("Start Loading gallaries.\n")
	galleries = map[int]*Gallery{}
	tags = map[TagK]map[TagV][]int{}
	tagBarrier := sync.WaitGroup{}
	tagBarrier.Add(1)
	go tagJson(jGalleriesCh, &tagBarrier)

	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer func() { _ = f.Close() }()
	//skip first {
	if _, e := f.Seek(1, 0); e != nil {
		return e
	}

	b := bufio.NewScanner(f)
	b.Split(scanJson)

	count := 0
	for b.Scan() {
		rawJsonCh <- b.Text()
		count++
		if count%10000 == 0 {
			logger.Printf("Parsed %d gallaries.\n", count)
		}
	}
	close(rawJsonCh)
	feedBarrier.Wait()
	logger.Printf("Parsed %d gallaries.\n", count)
	logger.Printf("End Parsing json.\n")
	close(jGalleriesCh)
	tagBarrier.Wait()
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
		case "male", "female", "misc", "language":
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

func scanJson(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	inColumn := false
	level, start, end := 0, 0, 0
	for l, s := range data {
		if s == '"' {
			if l == 0 || data[l-1] != '\\' {
				inColumn = !inColumn
			} else {
				logger.Printf("%d, %s\n", l, data)
			}
		}
		if s == '{' && !inColumn {
			if start == 0 {
				start = l
			}
			level++
		}
		if s == '}' && !inColumn {
			level--
			if level <= 0 {
				end = l
				break
			}
		}
	}
	if end > start {
		return end + 1, data[start : end+1], nil
		// We have a full newline-terminated line.
	} //logger.Println(start,end,level,len(data))

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), nil, nil
	}
	// Request more data.
	return 0, nil, nil
}

func feedJson(feedCh chan string, jsonCh chan JGallery, barrier *sync.WaitGroup) {
	for b := range feedCh {
		j := JGallery{}
		e := json.Unmarshal([]byte(b), &j)
		if e != nil {
			logger.Printf("Json unmarshal error: %s", e)
			continue
		}
		jsonCh <- j
	}
	barrier.Done()
}

func tagJson(jsonCh chan JGallery, barrier *sync.WaitGroup) {
	for j := range jsonCh {
		handleJGallery(j)
	}
	barrier.Done()
}
