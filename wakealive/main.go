package main

import (
	"fmt"
	"os"
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
	proc, err := lib.Find("dropbox")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%d %s", proc.Pid(), proc.Binary())
}
