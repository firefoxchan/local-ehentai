package ehloader

import "testing"

func TestParseTitle(t *testing.T) {
	titles := map[string]map[string]string{
		`123.(C99) [The.GR OU P (ArTi S T)] tit L E (P a.ro,dy) [翻訳] [translator].zip`: {
			rIdxTitleConvention:  "c99",
			rIdxTitleGroup:       "the.gr ou p",
			rIdxTitleArtist:      "arti s t",
			rIdxTitleTitle:       "tit l e",
			rIdxTitleParody:      "p a.ro,dy",
			rIdxTitleTranslation: "翻訳",
		},
		`[artist] title (parody) [翻訳]`: {
			rIdxTitleConvention:  "",
			rIdxTitleGroup:       "",
			rIdxTitleArtist:      "artist",
			rIdxTitleTitle:       "title",
			rIdxTitleParody:      "parody",
			rIdxTitleTranslation: "翻訳",
		},
		`[artist] title [翻訳]`: {
			rIdxTitleConvention:  "",
			rIdxTitleGroup:       "",
			rIdxTitleArtist:      "artist",
			rIdxTitleTitle:       "title",
			rIdxTitleParody:      "",
			rIdxTitleTranslation: "翻訳",
		},
		`[artist] title (parody)`: {
			rIdxTitleConvention:  "",
			rIdxTitleGroup:       "",
			rIdxTitleArtist:      "artist",
			rIdxTitleTitle:       "title",
			rIdxTitleParody:      "parody",
			rIdxTitleTranslation: "",
		},
		`title (parody)`: {
			rIdxTitleConvention:  "",
			rIdxTitleGroup:       "",
			rIdxTitleArtist:      "",
			rIdxTitleTitle:       "title",
			rIdxTitleParody:      "parody",
			rIdxTitleTranslation: "",
		},
		`title`: {
			rIdxTitleConvention:  "",
			rIdxTitleGroup:       "",
			rIdxTitleArtist:      "",
			rIdxTitleTitle:       "title",
			rIdxTitleParody:      "",
			rIdxTitleTranslation: "",
		},
	}
	for title, target := range titles {
		ret := parseTitle(title)
		mismatched := make([]string, 0, 5)
		for _, rIdx := range rIdxTitleAll {
			if ret[rIdx] != target[rIdx] {
				mismatched = append(mismatched, rIdx)
			}
		}
		if len(mismatched) != 0 {
			t.Errorf("Mismatch, %v, %s\n  Target: %+v\n     Ret: %+v", mismatched, title, target, ret)
		} else {
			t.Logf("Match, %s, %+v", title, ret)
		}
	}
}
