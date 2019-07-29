package websvr

import "testing"

func TestGenRatingStyle(t *testing.T) {
	cases := map[float32]string{
		0.0: "background-position:-80px -1px",
		0.1: "background-position:-80px -1px",
		0.4: "background-position:-64px -21px",
		0.5: "background-position:-64px -21px",
		0.6: "background-position:-64px -21px",
		1.0: "background-position:-64px -1px",
		1.5: "background-position:-48px -21px",
		2.0: "background-position:-48px -1px",
		2.5: "background-position:-32px -21px",
		3.0: "background-position:-32px -1px",
		3.5: "background-position:-16px -21px",
		4.0: "background-position:-16px -1px",
		4.5: "background-position:0px -21px",
		5.0: "background-position:0px -1px",
	}
	for rating, target := range cases {
		ret := genRatingStyle(rating)
		if ret != target {
			t.Errorf("Mismatch, Rating %f, Target %s, Ret %s", rating, target, ret)
		}
	}
}
