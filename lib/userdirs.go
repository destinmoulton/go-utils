package lib

import (
	"os"
	"path/filepath"
)

type XDGDirs struct {
	config string
	cache  string
}

var UserDirs XDGDirs

func init() {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		UserDirs.config = os.Getenv("XDG_CONFIG_HOME")
	} else {
		UserDirs.config = filepath.Join(os.Getenv("HOME"), ".config")
	}
	if os.Getenv("XDG_CACHE_HOME") != "" {
		UserDirs.cache = os.Getenv("XDG_CACHE_HOME")
	} else {
		UserDirs.cache = filepath.Join(os.Getenv("HOME"), ".cache")
	}
}

func (x *XDGDirs) Config() string {
	return x.config
}

func (x *XDGDirs) Cache() string {
	return x.cache
}
