// Package store is a dead simple configuration manager for Go applications.
package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

var applicationName string

// SetApplicationName defines a unique application handle for file system.
//
// By default, Store puts all your config data to %APPDATA%/<appname> on Windows
// and to $XDG_CONFIG_HOME or $HOME on *unix systems.
//
// Warning: Store would panic on any sensitive calls if it's not set.
func SetApplicationName(handle string) {
	applicationName = handle
}

// Load reads a configuration from `path` and puts it into `v` pointer.
//
// Path is a full filename, with extension. Since Store currently support
// TOML and JSON only, passing others would result in a corresponding error.
//
// If `path` doesn't exist, Load will create one and emptify `v` pointer by
// replacing it with a newly created object, derived from type of `v`.
func Load(path string, v interface{}) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	globalPath := buildPlatformPath(path)

	data, err := ioutil.ReadFile(globalPath)

	if err != nil {
		// There is a chance that file we are looking for
		// just doesn't exist. In this case we are supposed
		// to create an empty configuration file, based on v.
		empty := reflect.New(reflect.TypeOf(v))
		if innerErr := Save(path, &empty); innerErr != nil {
			// Must be smth with file system... returning error from read.
			return err
		}

		v = empty

		return nil
	}

	contents := string(data)
	if strings.HasSuffix(path, ".toml") {
		if _, err := toml.Decode(contents, v); err != nil {
			return err
		}
	} else if strings.HasSuffix(path, ".json") {
		if err := json.Unmarshal(data, v); err != nil {
			return err
		}
	} else {
		return &stringError{"unknown configuration format"}
	}

	return nil
}

// Save puts a configuration from `v` pointer into a file `path`.
//
// Path is a full filename, with extension. Since Store currently support
// TOML and JSON only, passing others would result in a corresponding error.
func Save(path string, v interface{}) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	var data []byte

	if strings.HasSuffix(path, ".toml") {
		var b bytes.Buffer

		encoder := toml.NewEncoder(&b)
		if err := encoder.Encode(v); err != nil {
			return nil
		}

		data = b.Bytes()
	} else if strings.HasSuffix(path, ".json") {
		fileData, err := json.Marshal(v)

		if err != nil {
			return err
		}

		data = fileData
	} else {
		return &stringError{"unknown configuration format"}
	}

	globalPath := buildPlatformPath(path)
	if err := os.MkdirAll(filepath.Dir(globalPath), os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(globalPath, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// buildPlatformPath builds a platform-dependent path for relative path given.
func buildPlatformPath(path string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\%s\\%s", os.Getenv("APPDATA"),
			applicationName,
			path)
	}

	var unixConfigDir string
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		unixConfigDir = xdg
	} else {
		unixConfigDir = os.Getenv("HOME") + "/.config"
	}

	return fmt.Sprintf("%s/%s/%s", unixConfigDir,
		applicationName,
		path)
}
