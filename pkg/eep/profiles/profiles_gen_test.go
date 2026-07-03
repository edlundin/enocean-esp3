package profiles

import "testing"

func TestGeneratedRegistryLoaded(t *testing.T) {
	if len(Registry) == 0 { t.Fatal("empty registry") }
	if _, ok := Registry["D5-00-01"]; !ok { t.Fatal("missing manual profile metadata") }
}
