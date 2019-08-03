package ehloader

import "testing"

func TestParseTitle(t *testing.T) {
	titles := map[string]Title{
		`(C99) [The.GR OU P (ArTi S T)] tit L E (P a.ro,dy) [翻訳]`: {
			Convention:  "c99",
			Group:       "the.gr ou p",
			Artist:      "arti s t",
			Title:       "tit l e",
			Parody:      "p a.ro,dy",
			Translation: "翻訳",
		},
		`[artist] title (parody) [翻訳]`: {
			Artist:      "artist",
			Title:       "title",
			Parody:      "parody",
			Translation: "翻訳",
		},
		`[artist] title [翻訳]`: {
			Artist:      "artist",
			Title:       "title",
			Translation: "翻訳",
		},
		`[artist] title (parody)`: {
			Artist:      "artist",
			Title:       "title",
			Parody:      "parody",
		},
		`title (parody)`: {
			Title:       "title",
			Parody:      "parody",
		},
		`title`: {
			Title:       "title",
		},
	}
	for title, target := range titles {
		ret := parseTitle(title)
		mismatched := make([]string, 0, 5)
		if ret.Convention != target.Convention {
			mismatched = append(mismatched, "Convention")
		}
		if ret.Group != target.Group {
			mismatched = append(mismatched, "Group")
		}
		if ret.Artist != target.Artist {
			mismatched = append(mismatched, "Artist")
		}
		if ret.Title != target.Title {
			mismatched = append(mismatched, "Title")
		}
		if ret.Parody != target.Parody {
			mismatched = append(mismatched, "Parody")
		}
		if ret.Translation != target.Translation {
			mismatched = append(mismatched, "Translation")
		}
		if len(mismatched) != 0 {
			t.Errorf("Mismatch, %v, %s\n  Target: %+v\n     Ret: %+v", mismatched, title, target, ret)
		} else {
			t.Logf("Match, %s\n  Ret: %+v", title, ret)
		}
	}
}
