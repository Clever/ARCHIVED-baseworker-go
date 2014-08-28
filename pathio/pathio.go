/*
pathio is a package that allows writing to and reading from different types of paths transparently.
It supports three types of paths:
 1. Empty Paths (stdin / stdout)
 2. Local file paths
 3. S3 File Paths (s3://bucket/object)

Note that using s3 paths requires setting three environment variables
 1. AWS_SECRET_ACCESS_KEY
 2. AWS_ACCESS_KEY_ID
 3. AWS_REGION
*/
package pathio

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

// ReaderForPath returns an io.Reader for the specified path. The path can either be empty (stdin),
// a local file path, or an S3 path.
func ReaderForPath(path string) (io.Reader, error) {
	if len(path) == 0 {
		return os.Stdin, nil
	} else if strings.HasPrefix(path, "s3://") {
		return s3FileReader(path)
	} else {
		// Local file path
		return os.Open(path)
	}
}

// WriteToPath writes a byte array to the specified path. The path can be either empty (stdout),
// a local file path, or an S3 path.
func WriteToPath(path string, input []byte) error {
	return WriteReaderToPath(path, bytes.NewReader(input), int64(len(input)))
}

// WriteReaderToPath writes all the data read from the specified io.Reader to the
// output path. The path can either be empty (stdout), a local file path, or an S3 path.
func WriteReaderToPath(path string, input io.Reader, length int64) error {
	if len(path) == 0 {
		_, err := io.Copy(os.Stdout, input)
		return err
	} else if strings.HasPrefix(path, "s3://") {
		return writeToS3(path, input, length)
	} else {
		return writeToLocalFile(path, input)
	}

}

func writeToLocalFile(path string, input io.Reader) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, input)
	return err

}

// s3FileReader converts an S3Path into an io.Reader
func s3FileReader(path string) (io.Reader, error) {
	bucket, s3path, err := getS3BucketAndPath(path)
	if err != nil {
		return nil, err
	}
	log.Printf("Getting from s3: %s", s3path)
	s3data, err := bucket.Get(s3path)
	if err != nil {
		log.Fatalf("Error downloading s3path: ", err)
		return nil, err
	}
	return bytes.NewReader(s3data), nil
}

func writeToS3(path string, input io.Reader, length int64) error {
	bucket, objectPath, err := getS3BucketAndPath(path)
	if err != nil {
		return err
	}
	log.Printf("Putting to s3: %s", path)
	return bucket.PutReader(objectPath, input, length, "text/plain", s3.Private)
}

// getS3BucketAndObject takes in a full s3path (s3://bucket/object) and returns a bucket,
// object name, error tuple. It assumes that AWS environment variables are set.
func getS3BucketAndPath(path string) (*s3.Bucket, string, error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatalf("AWS environment variables not set")
		return nil, "", err
	}

	region, err := region(os.Getenv("AWS_REGION"))
	// This is a HACK, but the S3 library we use doesn't support redirections from Amazon, so when
	// we make a request to https://s3-us-west-1.amazonaws.com and Amazon returns a 301 redirecting
	// to https://s3.amazonaws.com the library blows up.
	region.S3Endpoint = "https://s3.amazonaws.com"
	if err != nil {
		return nil, "", err
	}
	s := s3.New(auth, region)
	bucketName, s3path, err := parseS3Path(path)
	if err != nil {
		return nil, "", err
	}
	bucket := s.Bucket(bucketName)
	return bucket, s3path, err
}

// parseS3path parses an S3 path (s3://bucket/object) and returns a bucket, objectPath, error tuple
func parseS3Path(path string) (string, string, error) {
	// S3 path names are of the form s3://bucket/path
	stringsArray := strings.SplitAfterN(path, "/", 4)
	if len(stringsArray) < 4 {
		return "", "", fmt.Errorf("Invalid s3 path %s", path)
	}
	bucketName := stringsArray[2]
	// Strip off the slash
	bucketName = bucketName[0 : len(bucketName)-1]
	objectPath := stringsArray[3]
	return bucketName, objectPath, nil
}

// getRegion converts a region name into an aws.Region object
func region(regionString string) (aws.Region, error) {
	for name, region := range aws.Regions {
		if strings.ToLower(name) == strings.ToLower(regionString) {
			return region, nil
		}
	}
	return aws.Region{}, fmt.Errorf("Unknown region %s: ", regionString)
}
