package workerwrapper

import (
	"bufio"
	"bytes"
	"container/ring"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/Clever/baseworker-go"
	"github.com/Clever/baseworker-go/argsparser"
)

// TaskConfig defines the configuration for the task.
type TaskConfig struct {
	FunctionName, FunctionCmd string
	WarningLines              int
}

// Process runs the Gearman job by running the configured task.
func (conf TaskConfig) Process(job baseworker.Job) ([]byte, error) {

	log.Printf("Running job: %s:%s", conf.FunctionName, job.UniqueId())
	defer func() {
		// If we panicked then set the panic message as a warning. Gearman-go will
		// handle marking this job as failed.
		if r := recover(); r != nil {
			err := r.(error)
			job.SendWarning([]byte(err.Error()))
		}
	}()
	args, err := argsparser.ParseArgs(string(job.Data()))
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(conf.FunctionCmd, args...)
	var stderrbuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrbuf)
	var stdoutbuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutbuf)
	io.WriteString(cmd.Stderr, "input\n")
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	stderrToWarnings(stderrbuf, job, conf.WarningLines)
	job.SendData(stdoutbuf.Bytes())

	log.Printf("Finished job: %s:%s", conf.FunctionName, job.UniqueId())
	return nil, nil
}

func stderrToWarnings(buffer bytes.Buffer, job baseworker.Job, warningLines int) error {
	scanner := bufio.NewScanner(&buffer)
	lastStderrLines := ring.New(warningLines)
	for scanner.Scan() {
		lastStderrLines = lastStderrLines.Next()
		lastStderrLines.Value = scanner.Bytes()
	}
	for i := 0; i < lastStderrLines.Len(); i++ {
		lastStderrLines = lastStderrLines.Next()
		if lastStderrLines.Value != nil {
			job.SendWarning(lastStderrLines.Value.([]byte))
		}
	}
	return scanner.Err()
}

// processAndForwardStderr forwards the stderr from the worker process to the stderr of this process
// and also keeps track of the last X stderr lines. These last stderr lines are set to the warnings
// field in the Gearman job.
func (conf TaskConfig) processAndForwardStderr(job baseworker.Job, stderr io.Reader) error {
	scanner := bufio.NewScanner(stderr)
	lastStderrLines := ring.New(conf.WarningLines)
	for scanner.Scan() {
		lastStderrLines = lastStderrLines.Next()
		lastStderrLines.Value = scanner.Bytes()
	}
	for i := 0; i < lastStderrLines.Len(); i++ {
		lastStderrLines = lastStderrLines.Next()
		if lastStderrLines.Value != nil {
			job.SendWarning(lastStderrLines.Value.([]byte))
		}
	}
	return scanner.Err()
}
