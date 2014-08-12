package workerwrapper

import (
	"bufio"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/Clever/baseworker-go"
)

// WorkerConfig defines the configuration for the worker process this program wraps.
type WorkerConfig struct {
	JobName      string
	JobCmd       string
	WarningLines int
}

// Process runs the Gearman job by running the configured worker process.
func (conf WorkerConfig) Process(job baseworker.Job) (data []byte, err error) {

	log.Printf("Running job: %s:%s", conf.JobName, job.UniqueId())
	defer func() {
		// If we panicked then set the panic message as a warning. Gearman-go will
		// handle marking this job as failed.
		if r := recover(); r != nil {
			err = r.(error)
			job.SendWarning([]byte(err.Error()))
		}
	}()
	input := base64.StdEncoding.EncodeToString(job.Data())
	cmd := exec.Command(conf.JobCmd, input)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	err = conf.processAndForwardStderr(job, stderr)
	if err != nil {
		return nil, err
	}
	response, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	// This returns an error if the subprocess returns a non-zero exit code.
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	log.Printf("Finished job: %s:%s", conf.JobName, job.UniqueId())
	return response, nil
}

// processAndForwardStderr forwards the stderr from the worker process to the stderr of this process
// and also keeps track of the last X stderr lines. These last stderr lines are set to the warnings
// field in the Gearman job.
func (conf WorkerConfig) processAndForwardStderr(job baseworker.Job, stderr io.Reader) (err error) {
	scanner := bufio.NewScanner(stderr)
	lastStderrLines := make([][]byte, conf.WarningLines)
	for scanner.Scan() {
		os.Stderr.Write(scanner.Bytes())
		// If we already have enough warning lines then remove the first one (the oldest)
		// one and add the new one to the end
		if len(lastStderrLines) == conf.WarningLines {
			lastStderrLines = lastStderrLines[1:len(lastStderrLines)]
		}
		lastStderrLines = append(lastStderrLines, scanner.Bytes())
	}
	for index := range lastStderrLines {
		job.SendWarning(lastStderrLines[index])
	}
	return scanner.Err()
}
