package main

import (
	"fmt"
	"go-utils/lib"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

const NUM_LOGS_TO_KEEP = 25

type BackupObj struct {
	Name       string
	Group      string
	Src        string
	Dest       string
	Exclusions []string
	IsSSH      bool `json:"is_ssh"`
}

func main() {
	lib.UserDirs.Init("backup")
	lib.Logger.SetupLogger(lib.UserDirs.Logs())
	fmt.Println("starting")

	cpath := filepath.Join(lib.UserDirs.Config(), "backup.json")
	var backupArray []BackupObj
	lib.Config.ParseJSONToBytes(cpath, &backupArray)

	for _, backup := range backupArray {
		fmt.Println(backup.Name)
		args := []string{}
		args = append(args, "-avr")
		for _, excl := range backup.Exclusions {
			args = append(args, fmt.Sprintf("--exclude %s", excl))
		}
		if backup.IsSSH {
			args = append(args, "-e ssh")
		}
		args = append(args, backup.Src)
		args = append(args, backup.Dest)
		cmd := exec.Command("/usr/bin/rsync", args...)

		// Get (or make) the rsync log dir
		dir := getLogDir(&backup)
		fmt.Println(cmd)
		removeOldLogFiles(dir)
	}
}

func getLogDir(b *BackupObj) string {
	dir := filepath.Join(lib.UserDirs.Logs(), b.Group, b.Name)
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("unable to create log directory %s: %v", dir, err)
		}
	}
	return dir
}

func getLogFilename(b *BackupObj) string {
	dt := time.Now()
	timestamp := dt.Format("D01_02_2006-T15_04_05")
	return fmt.Sprintf("%s-%s.log", b.Name, timestamp)
}
func createLogFile(b *BackupObj) {

}
func removeOldLogFiles(path string) {
	fmt.Printf("removing for %s\n", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("unable to get the files in the %s path: %v", path, err)
	}
	// Sort the files - newest first
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	for i, file := range files {
		if (i + 1) > NUM_LOGS_TO_KEEP {
			fullpath := filepath.Join(path, file.Name())
			err := os.Remove(fullpath)
			if err != nil {
				log.Errorf("unable to remove log file %s: %v", fullpath, err)
			}
			log.Infof("removed old log file: %s", fullpath)
		}
	}
}
