package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	pid     int
	cmdline string
	binary  string
	ppid    int
	state   rune
	pgrp    int
	sid     int
}

type Procs struct {
}

func (p *Process) Pid() int {
	return p.pid
}
func (p *Process) Binary() string {
	return p.binary
}
func (p *Process) Cmdline() string {
	return p.cmdline
}

var Processes Procs

// Find a process by string in either the binary (preferred) or cmdline
func (p *Procs) Find(needle string, ignorepid int) (*Process, error) {
	processes, err := p.getAllProcesses()
	if err != nil {
		return nil, err
	}

	for _, proc := range processes {
		if proc.pid == ignorepid {
			// don't pay attention to this pid
			continue
		}
		if strings.Contains(proc.binary, needle) {
			return proc, nil
		}
		if strings.Contains(proc.cmdline, needle) {
			return proc, nil
		}
	}
	return nil, nil
}

// Determine if a process is running
func (p *Procs) IsRunning(name string, ignorepid int) (bool, error) {
	proc, err := p.Find(name, ignorepid)
	if err != nil {
		return false, err
	}
	if proc == nil {
		return false, nil
	}
	return true, nil
}

func (p *Procs) getAllProcesses() ([]*Process, error) {

	pf, err := os.Open("/proc")

	if err != nil {
		return nil, fmt.Errorf("could not open /proc %v", err)
	}
	defer pf.Close()

	var processes []*Process
	for {
		folders, err := pf.Readdirnames(20)
		if err == io.EOF {
			break
		}

		for _, folder := range folders {

			// Processes are numeric; ignore the rest
			if folder[0] < '0' || folder[0] > '9' {
				continue
			}

			pidint64, err := strconv.ParseInt(folder, 10, 0)
			if err != nil {
				continue
			}

			pid := int(pidint64)

			proc := &Process{pid: pid}
			err = proc.collate()
			if err != nil {
				continue
			}
			processes = append(processes, proc)
		}
	}
	return processes, nil
}

// collate builds the Process data from fs /proc files
func (p *Process) collate() error {
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", p.pid)
	cmdBytes, err := ioutil.ReadFile(cmdlineFile)
	if err != nil {
		return err
	}
	cmdEnd := 0
	for i, b := range cmdBytes {
		if int(b) == 0 {
			cmdEnd = i
			break
		}
		cmdEnd = i
	}

	cmdata := string(cmdBytes[:cmdEnd])
	p.cmdline = cmdata

	statPath := fmt.Sprintf("/proc/%d/stat", p.pid)
	sdBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}

	// First, parse out the image name
	data := string(sdBytes)
	binStart := strings.IndexRune(data, '(') + 1
	binEnd := strings.IndexRune(data[binStart:], ')')
	p.binary = data[binStart : binStart+binEnd]

	// Move past the image name and start parsing the rest
	data = data[binStart+binEnd+2:]
	_, err = fmt.Sscanf(data,
		"%c %d %d %d",
		&p.state,
		&p.ppid,
		&p.pgrp,
		&p.sid)

	return err
}
