package ehloader

import "testing"

func TestIntersect(t *testing.T) {
	matches := [][]int{
		{2, 3, 5, 10, 100, 200, 1000},
		{1, 2, 3, 5, 10, 200},
		{1, 5, 11, 200},
	}
	target := []int{5, 200}
	checkMatch(t, target, intersect(matches))
}

func TestUnion(t *testing.T) {
	matches := [][]int{
		{1, 5, 5, 11, 200},
		{},
		{2, 3, 5, 10, 100, 200, 1000},
		{1, 2, 3, 5, 10, 200},
	}
	target := []int{1, 2, 3, 5, 10, 11, 100, 200, 1000}
	checkMatch(t, target, union(matches))
}

func TestRSlice(t *testing.T) {
	match := []int{
		10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	}
	checkMatch(t, []int{0, 1, 2}, rSlice(match, 0, 3))
	checkMatch(t, []int{2, 3, 4}, rSlice(match, 2, 3))
	checkMatch(t, []int{8, 9, 10}, rSlice(match, 8, 3))
	checkMatch(t, []int{9, 10}, rSlice(match, 9, 3))
}

func checkMatch(t *testing.T, target, ret []int) {
	if len(ret) != len(target) {
		t.Errorf("Mismatch, Target: %v, Ret: %v", target, ret)
		return
	}
	mismatch := false
	for i, v := range ret {
		if v != target[i] {
			mismatch = true
		}
	}
	if mismatch {
		t.Errorf("Mismatch, Target: %v, Ret: %v", target, ret)
		return
	}
	t.Logf("Match, Target: %v, Ret: %v", target, ret)
}
