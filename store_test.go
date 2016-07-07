package store

import (
	"os"
	"testing"
)

type Cat struct {
	Name string
	Big  bool
}

type Settings struct {
	Age          int
	Cats         []Cat
	RandomString string
}

func equal(a, b Settings) bool {
	if a.Age != b.Age {
		return false
	}

	if a.RandomString != b.RandomString {
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
	Init("store_test")

	settings := Settings{
		Age: 42,
		Cats: []Cat{
			Cat{"Rudolph", true},
			Cat{"Patrick", false},
			Cat{"Jeremy", true},
		},
		RandomString: "gophers are gonna conquer the world",
	}

	settingsFile := "path/to/preferences.toml"

	err := Save(settingsFile, &settings)
	if err != nil {
		t.Fatalf("failed to save preferences: %s\n", err)
		return
	}

	defer os.Remove(buildPlatformPath(settingsFile))

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
