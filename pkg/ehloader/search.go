package ehloader

import (
	"strings"
)

func Search(searchTags map[TagK]map[TagV]struct{}, offset, limit int) []*Gallery {
	indexMu.RLock()
	defer indexMu.RUnlock()
	matches := make([][]int, 0)
	for k, vs := range searchTags {
		switch k {
		case TagKAll:
			localMatches := make([][]int, 0)
			for k := range tags {
				for v := range vs {
					if match, ok := matchTagKV(k, v, matchModeLike); ok {
						localMatches = append(localMatches, match)
					}
				}
			}
			matches = append(matches, union(localMatches))
		default:
			if _, ok := tags[k]; ok {
				for v := range vs {
					if match, ok := matchTagKV(k, v, matchModeLike); ok {
						matches = append(matches, match)
					}
				}
			}
		}
	}
	match := intersect(matches)
	gs := make([]*Gallery, 0, limit)
	for _, id := range rSlice(match, offset, limit) {
		if g, ok := galleries[id]; ok {
			gs = append(gs, g)
		}
	}
	return gs
}

func rSlice(match []int, offset int, limit int) []int {
	from := len(match) - offset - limit
	to := from + limit
	if to > len(match) {
		to = len(match)
	}
	if from < 0 {
		from = 0
	}
	if to < 0 {
		return []int{}
	}
	logger.Printf("RSlice: [%d:%d] @ %d", from, to, len(match))
	newLimit := to - from
	sliced := make([]int, newLimit)
	for i, v := range match[from:to] {
		sliced[newLimit-i-1] = v
	}
	return sliced
}

type matchMode int

const (
	matchModeLike matchMode = 0
	matchModeEq   matchMode = 1
)

func matchTagKV(matchK, matchV string, mode matchMode) ([]int, bool) {
	matchK = strings.ToLower(matchK)
	matchV = strings.ToLower(matchV)
	if _, ok := tags[matchK]; !ok {
		logger.Printf("Match: %s, %s, 0", matchK, matchV)
		return nil, false
	}
	matches := make([][]int, 0)
	switch mode {
	case matchModeLike:
		for v, match := range tags[matchK] {
			if strings.Contains(v, matchV) {
				matches = append(matches, match)
			}
		}
	case matchModeEq:
		if match, ok := tags[matchK][matchV]; ok {
			matches = append(matches, match)
		}
	}
	match := union(matches)
	logger.Printf("Match: %d, %s, %s, %d", mode, matchK, matchV, len(match))
	return match, true
}

func union(matches [][]int) []int {
	ret := make([]int, 0)
	if len(matches) == 0 {
		return ret
	}
	cursors := make([]int, len(matches))
	for {
		allPass := true
		smallest := 0
		for i, cursor := range cursors {
			if cursor != len(matches[i]) {
				allPass = false
				smallest = matches[i][cursor]
				break
			}
		}
		if allPass {
			break
		}
		for i, match := range matches {
			if cursors[i] >= len(match) {
				continue
			}
			if match[cursors[i]] < smallest {
				smallest = match[cursors[i]]
			}
		}
		ret = append(ret, smallest)
		for i, match := range matches {
			for cursors[i] < len(match) && match[cursors[i]] == smallest {
				cursors[i] += 1
			}
		}
	}
	return ret
}

func intersect(matches [][]int) []int {
	ret := make([]int, 0)
	if len(matches) == 0 {
		return ret
	}
	minI := 0
	for i, match := range matches {
		if len(match) < minI {
			minI = i
		}
	}
	cursors := make([]int, len(matches))
	for _, item := range matches[minI] {
		founds := make([]bool, len(matches))
		for i, match := range matches {
			founds[i] = false
		searchLoop:
			for {
				if cursors[i] >= len(match) {
					break
				}
				switch {
				case match[cursors[i]] < item:
					// find next
					cursors[i] += 1
				case match[cursors[i]] == item:
					// found
					founds[i] = true
					// find next
					cursors[i] += 1
					break searchLoop
				case match[cursors[i]] > item:
					// exceed
					// do nothing
					break searchLoop
				}
			}
		}
		found := true
		for _, f := range founds {
			found = found && f
		}
		if found {
			ret = append(ret, item)
		}
	}
	return ret
}
