package profiles

import "testing"

// TestGeneratedRegistryLoaded verifies GeneratedRegistryLoaded behavior.
func TestGeneratedRegistryLoaded(t *testing.T) {
	if len(Registry) == 0 {
		t.Fatal("empty registry")
	}
	if _, ok := Registry["D5-00-01"]; !ok {
		t.Fatal("missing manual profile metadata")
	}
}

// TestGeneratedSpecialRanges verifies GeneratedSpecialRanges behavior.
func TestGeneratedSpecialRanges(t *testing.T) {
	for _, tc := range []struct {
		profile, shortcut string
		rawMax, sentinel  int
	}{{"A5-20-10", "CVAR", 100, 255}, {"D2-05-00", "POS", 100, 127}, {"D2-05-00", "ANG", 100, 127}} {
		found := false
		for _, f := range Registry[tc.profile].Fields {
			if f.Shortcut == tc.shortcut && f.RawMax == tc.rawMax {
				_, found = f.Enum(uint64(tc.sentinel))
				break
			}
		}
		if !found {
			t.Fatalf("%s %s range or sentinel missing", tc.profile, tc.shortcut)
		}
	}
}
