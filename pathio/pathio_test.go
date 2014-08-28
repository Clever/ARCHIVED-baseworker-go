package pathio

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseS3Path(t *testing.T) {
	bucketName, s3path, err := parseS3Path("s3://clever-files/directory/path")
	assert.Nil(t, err)
	assert.Equal(t, bucketName, "clever-files")
	assert.Equal(t, s3path, "directory/path")

	bucketName, s3path, err = parseS3Path("s3://clever-files/directory")
	assert.Nil(t, err)
	assert.Equal(t, bucketName, "clever-files")
	assert.Equal(t, s3path, "directory")
}

func TestParseInvalidS3Path(t *testing.T) {
	_, _, err := parseS3Path("s3://")
	assert.NotNil(t, err)

	_, _, err = parseS3Path("s3://ag-ge")
	assert.NotNil(t, err)
}

func TestStdinReader(t *testing.T) {
	reader, err := ReaderForPath("")
	assert.Nil(t, err)
	assert.Equal(t, os.Stdin, reader)
}

func TestFileReader(t *testing.T) {
	// Create a temporary file and write some data to it
	file, err := ioutil.TempFile("/tmp", "pathioFileReaderTest")
	text := "fileReaderTest"
	assert.Nil(t, err)
	ioutil.WriteFile(file.Name(), []byte(text), 0644)

	reader, err := ReaderForPath(file.Name())
	assert.Nil(t, err)
	line, _, err := bufio.NewReader(reader).ReadLine()
	assert.Nil(t, err)
	assert.Equal(t, string(line), text)
}

func TestWriteToStdout(t *testing.T) {
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()
	file, err := ioutil.TempFile("/tmp", "writeToStdoutTest")
	assert.Nil(t, err)
	defer os.Remove(file.Name())
	os.Stdout = file
	WriteToPath("", []byte("teststdout"))
	line, err := ioutil.ReadFile(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, "teststdout", string(line))
}

func TestWriteToFilePath(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "writeToPathTest")
	assert.Nil(t, err)
	defer os.Remove(file.Name())

	assert.Nil(t, WriteToPath(file.Name(), []byte("testout")))
	output, err := ioutil.ReadFile(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, "testout", string(output))
}

func TestRegion(t *testing.T) {
	regionObj, err := region("us-west-1")
	assert.Nil(t, err)
	assert.Equal(t, regionObj.EC2Endpoint, "https://ec2.us-west-1.amazonaws.com")

	regionObj, err = region("BadRegion")
	assert.NotNil(t, err)
}
