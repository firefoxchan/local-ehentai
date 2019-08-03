package ehloader

import (
	"fmt"
	"github.com/firefoxchan/local-ehentai/pkg/cache"
	"sort"
	"strings"
	"time"
)

var searchQCache = cache.NewCache(time.Minute)

func SearchQ(orderBy string, offset, limit int, q Q) ([]*Gallery, int) {
	cacheKey := fmt.Sprintf("%s:%s", orderBy, q.Dump("", "", ""))
	var match []int
	if cached, ok := searchQCache.Get(cacheKey, 10*time.Minute); ok {
		match = cached.([]int)
		logger.Printf("SearchQ (cached): %d, %d, %s", offset, limit, cacheKey)
	} else {
		match = searchQ(q)
		switch orderBy {
		case OrderByPosted:
			sort.Sort(sortPosted(match))
		case OrderByRating:
			sort.Sort(sortRating(match))
		}
		logger.Printf("SearchQ: %d, %d, %s", offset, limit, cacheKey)
		searchQCache.Set(cacheKey, match, 10*time.Minute)
	}
	sliced := rSlice(match, offset, limit)
	gs := make([]*Gallery, len(sliced))
	for i, id := range sliced {
		if g, ok := galleries[id]; ok {
			gs[i] = g
		}
	}
	return gs, len(match)
}

func searchQ(q Q) []int {
	switch q.op {
	case QOpAnd:
		return searchQAnd(q.subQs)
	case QOpOr:
		return searchQOr(q.subQs)
	case QOpLike:
		return searchQLike(q.k, q.v)
	case QOpEq:
		return searchQEq(q.k, q.v)
	}
	return make([]int, 0, 0)
}

func searchQAnd(subQs []Q) []int {
	matches := make([][]int, len(subQs))
	for i, subQ := range subQs {
		matches[i] = searchQ(subQ)
	}
	return intersect(matches)
}

func searchQOr(subQs []Q) []int {
	matches := make([][]int, len(subQs))
	for i, subQ := range subQs {
		matches[i] = searchQ(subQ)
	}
	return union(matches)
}

func searchQMatch(k TagK, v TagV, mode matchMode) []int {
	switch k {
	case TagKAll:
		matches := make([][]int, 0)
		for k := range tags {
			switch strings.ToLower(k) {
			case TagKExpunged, TagKMinRating, TagKExists, TagKExistsIn, TagKGId:
				// special tags
				continue
			default:
				// pass
			}
			if match, ok := matchTagKV(k, v, mode); ok {
				matches = append(matches, match)
			}
		}
		return union(matches)
	default:
		if match, ok := matchTagKV(k, v, mode); ok {
			return match
		}
	}
	return make([]int, 0, 0)
}

func searchQLike(k TagK, v TagV) []int {
	return searchQMatch(k, v, matchModeLike)
}

func searchQEq(k TagK, v TagV) []int {
	return searchQMatch(k, v, matchModeEq)
}
