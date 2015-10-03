# Store
>Store is a dead simple configuration manager for Go applications.

[![GoDoc](https://godoc.org/github.com/tucnak/store?status.svg)](https://godoc.org/github.com/tucnak/store)

I didn't like existing configuration management solution, like [globalconf](https://github.com/rakyll/globalconf), [tachyon](https://github.com/vektra/tachyon) or [viper](https://github.com/spf13/viper). First two just don't feel right and viper, imo, a little overcomplicated—definitely too much for small projects.

Here is a hot example of Store in the wild:
```go
package main

import (
	"github.com/tucnak/store"
	"log"
	"os"
)

type Settings struct {
	Age int
	Cats []string
	Pi float64
}

func init() {
	// By default, Store puts all your config data to $HOME/.config/<appname>
	// on *nix systems and to %APPDATA%/<appname> on Windows.
	//
	// Warning: Store would panic on any sensitive calls if it's not set.
	store.SetApplicationName("joecockerfanclub")
}

func main() {
	var settings Settings
	err := store.Load("preferences.toml", &settings)
	if err != nil {
		log.Printf("failed to load preferences: %s\n", err)
		os.Exit(1)
	}

	// Some magic...
	settings.Age++

	err := store.Save("preferences.toml", &settings)
	if err != nil {
		log.Printf("failed to save preferences: %s\n", err)
		os.Exit(1)
	}
}
```

Store also supports a variety of handy features, find them out with [godoc](https://godoc.org/github.com/tucnak/store).