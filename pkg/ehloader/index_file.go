package ehloader

import (
	"bufio"
	"fmt"
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
	appendMatch := func(fn, path string, gid int) {
		if _, ok := matchedFiles[gid]; !ok {
			matchedFiles[gid] = make([]string, 0)
		}
		matchedFiles[gid] = append(matchedFiles[gid], path)
	}
	unmatchedFileList := make([]string, 0)
	unmatchedFileTitles := map[string]map[string]string{}
	if e := filepath.Walk(dirPath, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			logger.Printf("Walk files dir error: %s", e)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		fn := info.Name()
		if gid, ok := fnMap[fn]; ok {
			appendMatch(fn, path, gid)
			return nil
		}
		// trim ext
		fn = fn[0 : len(fn)-len(filepath.Ext(fn))]
		if gid, ok := fnMap[fn]; ok {
			appendMatch(fn, path, gid)
			return nil
		}
		title := parseTitle(fn)
		if title[rIdxTitleTitle] != "" {
			qsTitle := make([]Q, 0)
			qsTitleJpn := make([]Q, 0)
			for _, rIdx := range rIdxTitleAll {
				if rIdx == rIdxTitleTitle {
					continue
				}
				if title[rIdx] != "" {
					qsTitle = append(qsTitle, Eq(TagKRIdxTitlePrefix+rIdx, title[rIdx]))
					qsTitleJpn = append(qsTitleJpn, Eq(TagKRIdxTitleJpnPrefix+rIdx, title[rIdx]))
				}
			}
			if len(qsTitle) != 0 {
				q := Or(And(qsTitle...), And(qsTitleJpn...))
				matched := searchQ(q)
				if strings.Contains(path, "DL") {
					logger.Printf("%s, %v", q.Dump("", "", ""), matched)
				}
				isMatched := false
				for _, gid := range matched {
					g := galleries[gid]
					if indexFileMatchTitle(g.TitleExt[rIdxTitleTitle], title[rIdxTitleTitle]) {
						appendMatch(fn, path, gid)
						isMatched = true
					}
					if indexFileMatchTitle(g.TitleJpnExt[rIdxTitleTitle], title[rIdxTitleTitle]) {
						// match
						appendMatch(fn, path, gid)
						isMatched = true
					}
					// TODO: edit distance based title match
				}
				if isMatched {
					return nil
				}
			}
		} else {
			// TODO: edit distance based full match
		}
		unmatchedFileTitles[path] = title
		unmatchedFileList = append(unmatchedFileList, path)
		return nil
	}); e != nil {
		return nil, e
	}
	if len(unmatchedFileTitles) > 0 {
		lines := make([]string, 0)
		for _, path := range unmatchedFileList {
			lines = append(lines, fmt.Sprintf("%s -> %v", path, unmatchedFileTitles[path]))
		}
		logger.Printf("Unmatched files: (%d)\n  %s", len(unmatchedFileTitles), strings.Join(lines, "\n  "))
	}
	return matchedFiles, nil
}

func indexFileMatchTitle(source, fn string) bool {
	sources := strings.Split(source, "|")
	if len(sources) > 1 {
		sources = append(sources, source)
	}
	fn = indexWinPathReplacer.Replace(fn)
	fn = strings.Join(split(fn, " "), " ")
	for _, source := range sources {
		source = indexWinPathReplacer.Replace(source)
		source = strings.Join(split(source, " "), " ")
		if strings.Contains(source, fn) {
			return true
		}
		if strings.Contains(fn, source) {
			return true
		}
	}
	return false
}

func split(s string, sep string) []string {
	sl := strings.Split(s, sep)
	ret := make([]string, 0, len(sl))
	for _, s := range sl {
		s = strings.TrimSpace(s)
		if s != "" {
			ret = append(ret, s)
		}
	}
	return ret
}

func indexFileMap(mapPath string) (map[string]int, error) {
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
