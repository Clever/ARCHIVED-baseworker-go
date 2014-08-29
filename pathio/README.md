# pathio
--
    import "github.com/Clever/baseworker-go/pathio"

pathio is a package that allows writing to and reading from different types of
paths transparently. It supports three types of paths:

    1. Empty Paths (stdin / stdout)
    2. Local file paths
    3. S3 File Paths (s3://bucket/object)

Note that using s3 paths requires setting three environment variables

    1. AWS_SECRET_ACCESS_KEY
    2. AWS_ACCESS_KEY_ID
    3. AWS_REGION

## Usage

#### func  ReaderForPath

```go
func ReaderForPath(path string) (io.Reader, error)
```
ReaderForPath returns an io.Reader for the specified path. The path can either
be empty (stdin), a local file path, or an S3 path.

#### func  WriteReaderToPath

```go
func WriteReaderToPath(path string, input io.Reader, length int64) error
```
WriteReaderToPath writes all the data read from the specified io.Reader to the
output path. The path can either be empty (stdout), a local file path, or an S3
path.

#### func  WriteToPath

```go
func WriteToPath(path string, input []byte) error
```
WriteToPath writes a byte array to the specified path. The path can be either
empty (stdout), a local file path, or an S3 path.

## Testing

You can run the test cases by typing `make test` in the root of the repository

## Documentation

The documentation is automatically generated via [godocdown](https://github.com/robertkrimen/godocdown).

You can update it by typing `make docs` in the root of the repository

They're also viewable online at [![GoDoc](https://godoc.org/github.com/Clever/baseworker-go/pathio?status.png)](https://godoc.org/github.com/Clever/baseworker-go/pathio).
