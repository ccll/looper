package gat

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Run struct {
	Tags    string
	Verbose bool
}

func (run Run) RunAll() {
	run.goTest("./...")
}

func (run Run) RunOnChange(file string) {
	if isGoFile(file) {
		// TODO: optimization, skip if no test files exist
		packageDir := "./" + filepath.Dir(file) // watchDir = ./
		run.goTest(packageDir)
	}
}

func (run Run) goTest(test_files string) {
	args := []string{"test"}
	if len(run.Tags) > 0 {
		args = append(args, []string{"-tags", run.Tags}...)
	}
	if run.Verbose {
		args = append(args, "-v")
	}
	args = append(args, test_files)

	command := "go"

	if _, err := os.Stat("Godeps/Godeps.json"); err == nil {
		args = append([]string{"go"}, args...)
		command = "godep"
	}

	cmd := exec.Command(command, args...)
	// cmd.Dir watchDir = ./

	PrintCommand(cmd.Args) // includes "go"

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	PrintCommandOutput(out)

	RedGreen(cmd.ProcessState.Success())
	ShowDuration(cmd.ProcessState.UserTime())
}

func isGoFile(file string) bool {
	return filepath.Ext(file) == ".go"
}
