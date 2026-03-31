package command

import (
	"sort"
	"testing"
)

func TestByVersion(t *testing.T) {
	versions := []string{"1.2.3", "0.1.0", "1.10.2", "1.1.9"}
	expected := []string{"0.1.0", "1.1.9", "1.2.3", "1.10.2"}

	sort.Sort(ByVersion(versions))

	for i, v := range versions {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

func TestSetWellKnownValues(t *testing.T) {
	dict := make(map[string]string)
	dict = SetWellKnownValues(dict)

	keys := []string{"latesttag", "actualdate", "actualtime", "actualtimestamp"}
	for _, k := range keys {
		if _, ok := dict[k]; !ok {
			t.Errorf("expected key %s to be set", k)
		}
	}
}
