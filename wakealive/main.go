package main

import (
	"fmt"
	"os"
	"os/exec"
	"utiligo/lib"
)

func main() {
	// processes, err := lib.GetAllProcesses()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }

	// for _, proc := range *processes {
	//if proc.pid == 742 {
	// 	fmt.Printf("%d %s %s\n", proc.Pid(), proc.Binary(), proc.Cmdline())
	//	}
	// }
	fmt.Printf("temp dir %s\n", os.TempDir())

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
	isWithin := lib.ImageTools.IsImageWithin(small, filename)
	os.Remove(filename)
	fmt.Printf("Is %s within %s? %v\n", small, filename, isWithin)

	// Check if the process is running
	p := "dropbox"
	isRunning, err := lib.Processes.IsRunning(p)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Is %s running? %v\n", p, isRunning)

	if !isWithin && isRunning {
		// Reload i3
		fmt.Println("Restarting i3")
		cmd := exec.Command("/usr/bin/i3-msg", "restart")
		cmd.Run()
	}
}
