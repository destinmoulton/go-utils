package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

type DirData struct {
	dirs map[string]string
}

var UserDirs DirData

func (d *DirData) Init(appname string) {
	module := "go-utils"
	d.dirs = make(map[string]string)
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		d.dirs["config"] = filepath.Join(os.Getenv("XDG_CONFIG_HOME"), module, appname)
	} else {
		d.dirs["config"] = filepath.Join(os.Getenv("HOME"), ".config", module, appname)
	}
	if os.Getenv("XDG_CACHE_HOME") != "" {
		d.dirs["cache"] = filepath.Join(os.Getenv("XDG_CACHE_HOME"), module, appname)
	} else {
		d.dirs["cache"] = filepath.Join(os.Getenv("HOME"), ".cache", module, appname)
	}

	d.dirs["local"] = filepath.Join(os.Getenv("HOME"), ".local", module, appname)
	d.dirs["logs"] = filepath.Join(d.dirs["local"], "logs")

	d.makeUserDirs()
}

func (d *DirData) Config() string {
	return d.dirs["config"]
}

func (d *DirData) Cache() string {
	return d.dirs["cache"]
}

func (d *DirData) Local() string {
	return d.dirs["local"]
}

func (d *DirData) Logs() string {
	return d.dirs["logs"]
}

func (d *DirData) makeUserDirs() {
	for dir := range d.dirs {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			fmt.Printf("error creating userdir %s %v", dir, err)
			os.Exit(1)
		}
	}
}
