package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

type Dirs struct {
	config string
	cache  string
	local  string
	logs   string
}

var UserDirs = Dirs{}

func (d *Dirs) Init(appname string) {
	module := "go-utils"
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		d.config = filepath.Join(os.Getenv("XDG_CONFIG_HOME"), module, appname)
	} else {
		d.config = filepath.Join(os.Getenv("HOME"), ".config", module, appname)
	}
	if os.Getenv("XDG_CACHE_HOME") != "" {
		d.cache = filepath.Join(os.Getenv("XDG_CACHE_HOME"), module, appname)
	} else {
		d.cache = filepath.Join(os.Getenv("HOME"), ".cache", module, appname)
	}

	d.local = filepath.Join(os.Getenv("HOME"), ".local", module, appname)
	d.logs = filepath.Join(d.local, "logs")

	makeUserDir(d.config)
	makeUserDir(d.cache)
	makeUserDir(d.local)
	makeUserDir(d.logs)
}

func (x *Dirs) Config() string {
	return x.config
}

func (x *Dirs) Cache() string {
	return x.cache
}

func (x *Dirs) Local() string {
	return x.local
}

func (x *Dirs) Logs() string {
	return x.logs
}

func makeUserDir(path string) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Printf("error creating userdir %s %v", path, err)
		os.Exit(1)
	}
}
