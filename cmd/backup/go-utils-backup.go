package main

import (
	"flag"
	"fmt"
	"go-utils/lib"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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
	SSHPemFile string `json:"ssh_pem_file"`
}
type BackupArray []BackupObj
type BackupMap map[string]BackupObj
type Opts struct {
	delete bool
	jobs   string
}

func main() {
	opts := Opts{delete: false, jobs: "all"}
	flag.BoolVar(&opts.delete, "d", false, "run rsync with --delete-excluded")
	flag.StringVar(&opts.jobs, "j", "all", "run the selected backup job(s) from the json file (comma separate multiple)")
	flag.Parse()

	lib.UserDirs.Init("backup")
	lib.Logger.SetupLogger(lib.UserDirs.Logs(), "backup")

	cpath := filepath.Join(lib.UserDirs.Config(), "backup.json")
	var backupArray BackupArray
	lib.Config.ParseJSONToBytes(cpath, &backupArray)

	backupsToRun := parseJobsOption(opts, backupArray)
	slines := []string{}
	for _, backup := range backupsToRun {
		// Get (or make) the rsync log dir
		dir := getLogDir(&backup)
		logpath := getFullLogPath(&backup)
		blines, command := generateBackupSection(&backup, logpath, &opts)
		slines = append(slines, blines...)
		createLogFile(&backup, logpath, fmt.Sprintf("%s\n\n", command))
		removeOldLogFiles(dir)
	}

	bfile := "backup.sh"
	shfile := filepath.Join(lib.UserDirs.Config(), bfile)
	err := writeShellScript(slines, shfile)
	if err != nil {
		log.Fatalf("failed to generate backup file %s: %v", shfile, err)
	}
	log.Infof("generated backup file: %s", shfile)

	// build the full command
	cmd := exec.Command("/bin/bash", shfile)
	log.Infof("running backup `%s`: %s", shfile, cmd)
	out, err := cmd.CombinedOutput() // NOTE: Run as Start to keep script running in case of this crashing
	log.Infof("backup completed with output: %s", string(out))
	if err != nil {
		// Don't Fatal; capture the errors in the log
		log.Infof("errors: %v", err)
	}
}
func parseJobsOption(opts Opts, orig BackupArray) BackupMap {
	backupsMap := make(BackupMap)
	for _, backup := range orig {
		backupsMap[backup.Name] = backup
	}
	backupsToRun := make(BackupMap)
	if opts.jobs == "all" {
		backupsToRun = backupsMap
	} else {
		poss := strings.Split(opts.jobs, ",")
		for _, s := range poss {
			b, ok := backupsMap[s]
			if !ok {
				log.Warnf("the backup job `%s` is not in the json config file", s)
			} else {
				backupsToRun[s] = b
			}
		}
	}
	return backupsToRun
}
func writeShellScript(lines []string, filename string) error {
	shlines := []string{
		"#!/bin/bash", // Run it as a bash script
		"#",
		"# This file was generated by go-utils-backup",
		"#",
		"",
	}

	shlines = append(shlines, lines...)
	return ioutil.WriteFile(filename, []byte(strings.Join(shlines, "\n")+"\n"), 0774)
}

// Generate the lines of a backup
func generateBackupSection(b *BackupObj, logpath string, opts *Opts) ([]string, string) {

	lines := []string{
		"",
	}
	// Add args
	hasExcl := false
	var excls []string
	exclude := ""
	for _, ex := range b.Exclusions {
		hasExcl = true
		// Wrap them
		excls = append(excls, fmt.Sprintf("--exclude='%s'", ex))
	}
	if hasExcl {
		exclude = strings.Join(excls, " ")
	}

	delarg := ""
	if opts.delete {
		delarg = "--delete-excluded"
	}

	ssh := ""
	if b.SSHPemFile != "" {
		ssh = fmt.Sprintf("-e 'ssh -i %s'", b.SSHPemFile)
	}

	// Add a comment
	comment := fmt.Sprintf("# Backup %s", b.Name)
	lines = append(lines, comment)

	// Pipe the rsync output into end of the log file
	command := fmt.Sprintf("rsync -avr %s %s %s %s %s >> %s", exclude, delarg, ssh, b.Src, b.Dest, logpath)
	return append(lines, command), command
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

func getFullLogPath(b *BackupObj) string {
	return filepath.Join(getLogDir(b), getLogFilename(b))
}
func createLogFile(b *BackupObj, logpath string, cmdstr string) {
	s := []byte(cmdstr)
	err := ioutil.WriteFile(logpath, s, 0666)
	if err != nil {
		log.Fatalf("unable to create log file `%s`: %v", logpath, err)
	}
}
func removeOldLogFiles(path string) {
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
