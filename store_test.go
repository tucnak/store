package store

import (
	"os"
	"testing"
)

func init() {
	SetApplicationName("store_test")
}

type Settings struct {
	Age  int
	Cats []string
	Pi   float64
}

func equal(a, b Settings) bool {
	if a.Age != b.Age {
		return false
	}

	if a.Pi != b.Pi {
		return false
	}

	if len(a.Cats) != len(b.Cats) {
		return false
	}

	for i, cat := range a.Cats {
		if cat != b.Cats[i] {
			return false
		}
	}

	return true
}

func TestSaveLoad(t *testing.T) {
	settings := Settings{
		Age:  42,
		Cats: []string{"cat1", "cat2", "cat3"},
		Pi:   3.1415,
	}

	settingsFile := "preferences.toml"

	err := Save(settingsFile, &settings)
	if err != nil {
		t.Fatalf("failed to save preferences: %s\n", err)
		return
	}

	defer os.Remove(buildPath(settingsFile))

	var newSettings Settings

	err = Load(settingsFile, &newSettings)
	if err != nil {
		t.Fatalf("failed to load preferences: %s\n", err)
		return
	}

	if !equal(settings, newSettings) {
		t.Fatalf("broken")
	}
}
