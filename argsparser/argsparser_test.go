package argsparser

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, expected, actual interface{}, message string) {
	message += fmt.Sprintf("\nExpected: '%s', got: '%s'", expected, actual)
	assert(t, expected == actual, message)
}

func assert(t *testing.T, assertion bool, message string) {
	if !assertion {
		t.Fatal(message)
	}
}

func assertIfErr(t *testing.T, err error, message string) {
	if err != nil {
		assert(t, false, fmt.Sprintf("%s: %s", message, err.Error()))
	}
}

func TestParseArgs(t *testing.T) {
	argsArray, err := ParseArgs("\"arg with quotes\" secondArg thirdArg \"another with quotes\"")
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(argsArray) != 4 {
		t.Fatal("Args length = " + strconv.Itoa(len(argsArray)) + ", 4 expected")
	}
	message := "Unexpected parsed argument"
	assertEqual(t, "arg with quotes", argsArray[0], message)
	assertEqual(t, "secondArg", argsArray[1], message)
	assertEqual(t, "thirdArg", argsArray[2], message)
	assertEqual(t, "another with quotes", argsArray[3], message)
}

func TestLongArgs(t *testing.T) {
	rawData, err := ioutil.ReadFile("/Users/sparks/dev/baseworker-go/argsparser/long_args.txt")
	assertIfErr(t, err, "Could not parse test data file")
	commandline := strings.Replace(string(rawData), "\n", "", -1)
	expectedArgs := strings.Split(commandline, " ")
	parsedArgs, err := ParseArgs(commandline)
	assertIfErr(t, err, "Error during ParseArgs")
	assertEqual(t, len(expectedArgs), len(parsedArgs), "Wrong number of arguments")
	for i, _ := range expectedArgs {
		assertEqual(t, expectedArgs[i], parsedArgs[i], "Unexpected parsed argument"+strconv.Itoa(i))
	}
}
