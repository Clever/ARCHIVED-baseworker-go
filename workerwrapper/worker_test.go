package workerwrapper

import (
	"encoding/base64"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	mock "github.com/Clever/baseworker-go/mock"
)

// Helper function to get the response for a job that should be successful
func getSuccessResponse(payload string, cmd string, t *testing.T) string {
	mockJob := mock.CreateMockJob(payload)
	config := WorkerConfig{JobName: "name", JobCmd: cmd, WarningLines: 5}
	response, err := config.Process(mockJob)
	if err != nil {
		t.Fatal(err)
	}
	return string(response)
}

// Helper function to assert that two strings are equal
func checkStringsEqual(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Fatal("Actual response: " + actual + " does not match expected: " + expected)
	}
}

// Helper function to assert that two integers and equal
func checkIntsEqual(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Fatal("Actual response: " + strconv.Itoa(actual) + " does not match expected: " +
			strconv.Itoa(expected))
	}
}

func TestSuccessResponse(t *testing.T) {
	response := getSuccessResponse("IgnorePayload", "testscripts/success.sh", t)
	checkStringsEqual(t, "SuccessResponse\n", response)
}

func TestErrorOnNonZeroExitCode(t *testing.T) {
	mockJob := mock.CreateMockJob("IgnorePayload")
	config := WorkerConfig{JobName: "name", JobCmd: "testscripts/nonZeroExit.sh", WarningLines: 5}
	response, err := config.Process(mockJob)
	if response != nil {
		t.Fatal("Should be no response on errored job")
	}
	if err == nil {
		t.Fatal("Job should have failed")
	}
	checkStringsEqual(t, "exit status 2", err.Error())
}

func TestWorkerRecievesInputData(t *testing.T) {
	response := getSuccessResponse("InputString", "testscripts/echoInput.sh", t)
	decodedResponse, err := base64.StdEncoding.DecodeString(response)
	if err != nil {
		t.Fatal(err.Error())
	}
	decodedString := string(decodedResponse)
	checkStringsEqual(t, "InputString", decodedString)
}

func TestStderrForwardedToProcess(t *testing.T) {
	// This test creates a child process because we want to make sure that the stderr of the worker
	// process is forwarded to the child process correctly. If we don't create a child process we
	// end up checking our own process' stderr which is a pain.
	cmd := exec.Command("go", "run", "testscripts/test_stderr.go")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cmd.Start(); err != nil {
		t.Fatal(err.Error())
	}
	response, err := ioutil.ReadAll(stderr)
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cmd.Wait(); err != nil {
		t.Fatal(err.Error())
	}
	if !strings.Contains(string(response), "stderr") {
		t.Fatal("Missing expected stderr output: " + string(response))
	}
}

func TestStderrCapturedInWarnings(t *testing.T) {
	mockJob := mock.CreateMockJob("IngnorePayload")
	config := WorkerConfig{JobName: "name", JobCmd: "testscripts/logStderr.sh", WarningLines: 2}
	_, err := config.Process(mockJob)
	if err != nil {
		t.Fatal(err)
	}
	warnings := mockJob.Warnings()
	checkIntsEqual(t, 2, len(warnings))
	checkStringsEqual(t, string(warnings[0]), "stderr7")
	checkStringsEqual(t, string(warnings[1]), "stderr8")
}

func TestHandleStderrAndStdoutTogether(t *testing.T) {
	mockJob := mock.CreateMockJob("IngnorePayload")
	config := WorkerConfig{JobName: "name", JobCmd: "testscripts/logStdoutAndStderr.sh", WarningLines: 5}
	response, err := config.Process(mockJob)
	if err != nil {
		t.Fatal(err)
	}
	warnings := mockJob.Warnings()
	lastWarning := warnings[len(warnings)-1]
	checkStringsEqual(t, "stderr2", string(lastWarning))
	checkStringsEqual(t, "stdout1\nstdout2\n", string(response))
}
