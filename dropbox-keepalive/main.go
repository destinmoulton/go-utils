package main

import (
	"fmt"
	"go-utils/lib"
	"os"
	"os/exec"
)

func main() {
	lib.UserDirs.Init("dropbox-keepalive")
	fmt.Println(lib.UserDirs.Config())
	// Check if the process is running
	p := "dropbox"
	isRunning, err := lib.Processes.IsRunning(p)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !isRunning {
		lib.DBus.Msg("Dropbox Alive?", "Dropbox is not in process list. Starting...")
		cmd := exec.Command("/usr/bin/dropbox")
		cmd.Run()
		os.Exit(0)
	}

	// Get the system tray and save as png
	img, err := lib.ImageTools.SystrayShot(120)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	filename, err := lib.ImageTools.SaveAsTempPNG(img)
	if err != nil {
		fmt.Printf("unable to save temp png: %v\n", err)
		os.Exit(1)
	}
	small := "dropbox_icon.png"
	isIconInSystray := lib.ImageTools.IsImageWithin(small, filename)
	os.Remove(filename)
	//fmt.Printf("Is %s within %s? %v\n", small, filename, isIconInSystray)

	//fmt.Printf("Is %s running? %v\n", p, isRunning)

	if !isIconInSystray && isRunning {
		// Reload i3
		cmd := exec.Command("/usr/bin/i3-msg", "restart")
		cmd.Run()
		lib.DBus.Msg("Dropbox Alive?", "i3 restarted for the Dropbox icon.")
	}
}