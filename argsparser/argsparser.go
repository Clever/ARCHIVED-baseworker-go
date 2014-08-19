package argsparser

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

// ParseArgs converts the command line specified into a slice of the command line arguments.
func ParseArgs(commandline string) ([]string, error) {
	file, err := ioutil.TempFile("/tmp", "parseArgs")
	fileClosed := false
	defer func() {
		if !fileClosed {
			file.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	// This is a bit hacky, but we couldn't think of a better way to do it.
	// We create a bash script and in that file we run a bash command that parses the
	// command line arguments we wrote to the file. The bash script outputs each of the
	// parsed arguments to stdout, separated by \n. We parse the stdout and return
	// that to the caller.
	file.WriteString("#!/bin/bash\n")
	file.WriteString("bash -c 'while test ${#} -gt 0; do echo $1; shift; done;' _ " + commandline + "\n")

	if err := file.Chmod(0744); err != nil {
		return nil, err
	}
	file.Close()
	fileClosed = true
	cmd := exec.Command(file.Name())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	response, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err = cmd.Wait(); err != nil {
		return nil, err
	}
	argsArray := strings.Split(string(response), "\n")
	// Remove the last element of the argsArray because the output ends with an endline
	// and has an empty last element
	argsArray = argsArray[0 : len(argsArray)-1]
	return argsArray, nil
}
