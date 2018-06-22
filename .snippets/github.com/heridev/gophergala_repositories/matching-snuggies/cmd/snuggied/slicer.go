package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SlicerCmd struct {
	Bin    string
	Args   []string
	OutLog io.Writer
	ErrLog io.Writer
}

type Slicer interface {
	SlicerCmd() *SlicerCmd
}

func Run(s Slicer, kill <-chan error) error {
	scmd := s.SlicerCmd()
	log.Printf("slicing with %s %v", scmd.Bin, scmd.Args)
	cmd := exec.Command(scmd.Bin, scmd.Args...)
	cmd.Stdout = scmd.OutLog
	cmd.Stderr = scmd.ErrLog
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("%s: %v", scmd.Bin, err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	for {
		select {
		case err := <-done:
			return err
		case err := <-kill:
			log.Printf("killing process %v", cmd.Process.Pid)
			if errkill := cmd.Process.Kill(); errkill != nil {
				// we couldn't kill the process. don't exit the loop.
				log.Printf("kill: %v", errkill)
				continue
			}
			return err
		}
	}
}

func ReadPresetsDirSlic3r(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".ini" {
			name := filepath.Base(file.Name())
			name = strings.TrimSuffix(name, ".ini")
			m[name] = filepath.Join(dir, file.Name())
		}
	}
	return m, nil
}

type Slic3r struct {
	Bin        string
	ConfigPath string
	OutPath    string
	InPath     string
}

func (s *Slic3r) SlicerCmd() *SlicerCmd {
	bin := s.Bin
	if bin == "" {
		bin = "slic3r"
	}
	var args []string
	config := s.ConfigPath
	if config != "" {
		args = append(args, "--load", config)
	}
	out := s.OutPath
	if out != "" {
		args = append(args, "-o", out)
	}
	in := s.InPath
	if in != "" {
		args = append(args, in)
	}
	return &SlicerCmd{
		Bin:    bin,
		Args:   args,
		OutLog: os.Stderr,
		ErrLog: os.Stderr,
	}
}
