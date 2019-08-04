package ehloader

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func indexFiles(dirPath string, mapPath string) (map[int][]string, error) {
	// load map
	fnMap, e := indexFileMap(mapPath)
	if e != nil {
		return nil, e
	}
	matchedFiles := map[int][]string{}
	unmatchedFiles := make([]string, 0)
	if e := filepath.Walk(dirPath, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.IsDir() {
			return nil
		}
		fn := info.Name()
		if gid, ok := fnMap[fn]; ok {
			if _, ok := matchedFiles[gid]; !ok {
				matchedFiles[gid] = make([]string, 0)
			}
			matchedFiles[gid] = append(matchedFiles[gid], path)
			return nil
		}
		// TODO: edit distance based match
		unmatchedFiles = append(unmatchedFiles, path)
		return nil
	}); e != nil {
		return nil, e
	}
	logger.Printf("Unmatched files: (%d)\n  %s", len(unmatchedFiles), strings.Join(unmatchedFiles, "\n  "))
	return matchedFiles, nil
}

func indexFileMap (mapPath string) (map[string]int, error) {
	fnMap := make(map[string]int, 0)
	if mapPath != "" {
		f, e := os.OpenFile(mapPath, os.O_RDONLY, 0)
		if e != nil {
			logger.Printf("Unable to load filename map: %s", e)
			return nil, e
		}
		defer func() { _ = f.Close() }()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "#") {
				continue
			}
			pairs := strings.SplitN(line, "/", 2)
			if len(pairs) != 2 {
				logger.Printf("Unable to parse filename map line: %s", line)
				continue
			}
			fn := strings.TrimSpace(pairs[0])
			gid, e := strconv.ParseInt(strings.TrimSpace(pairs[1]), 10, 64)
			if e != nil {
				logger.Printf("Unable to parse filename map gid: %s, %s", e, line)
				continue
			}
			fnMap[fn] = int(gid)
		}
		if e := scanner.Err(); e != nil {
			return nil, e
		}
	}
	return fnMap, nil
}
