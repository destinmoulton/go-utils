package main

import (
	"go-utils/lib"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func main() {
	lib.UserDirs.Init("dropbox-keepalive")
	lib.Logger.SetupLogger(lib.UserDirs.Logs(), "dropbox-keepalive")
	// Check if the process is running
	p := "dropbox"
	isRunning, err := lib.Processes.IsRunning(p, os.Getpid())
	if err != nil {
		log.Fatalf("error scanning processes: %v", err)
	}

	if !isRunning {
		log.Info("dropbox not in process list: starting dropbox")
		lib.DBus.Msg("Dropbox Keepalive", "Dropbox is not in process list. Starting...")
		cmd := exec.Command("/usr/bin/dropbox")
		cmd.Run()
		os.Exit(0)
	}

	// Get the system tray and save as png
	img, err := lib.ImageTools.SystrayShot(120)
	if err != nil {
		log.Fatalf("unable to make systray snapshot: %v", err)
	}
	filename, err := lib.ImageTools.SaveAsTempPNG(img)
	if err != nil {
		log.Fatalf("unable to save temp png: %v\n", err)
	}
	small := filepath.Join(lib.UserDirs.Config(), "dropbox_icon.png")
	isIconInSystray := lib.ImageTools.IsImageWithin(small, filename)
	err = os.Remove(filename)
	if err != nil {
		log.Fatalf("failed to remove the temporary screenshot: %v", err)
	}
	//fmt.Printf("Is %s within %s? %v\n", small, filename, isIconInSystray)

	//fmt.Printf("Is %s running? %v\n", p, isRunning)

	if !isIconInSystray && isRunning {
		// Reload i3
		log.Info("restarted i3: dropbox icon not in systray.")
		cmd := exec.Command("/usr/bin/i3-msg", "restart")
		cmd.Run()
		lib.DBus.Msg("Dropbox Keepalive", "i3 restarted for the Dropbox systray icon")
	}
}
