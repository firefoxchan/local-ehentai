package ehloader

import "sort"

const (
	OrderByGId    = ""
	OrderByPosted = "posted"
	OrderByRating = "rating"
)

type sortPosted []int

var _ sort.Interface = sortPosted{}

func (s sortPosted) Len() int      { return len(s) }
func (s sortPosted) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortPosted) Less(i, j int) bool {
	return galleries[s[i]].Posted.Before(galleries[s[j]].Posted)
}

type sortRating []int

var _ sort.Interface = sortRating{}

func (s sortRating) Len() int      { return len(s) }
func (s sortRating) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortRating) Less(i, j int) bool {
	return galleries[s[i]].Rating < galleries[s[j]].Rating
}
