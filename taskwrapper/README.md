# workerwrapper
--
    import "github.com/Clever/baseworker-go/workerwrapper"


## Usage

#### type TaskConfig

```go
type TaskConfig struct {
	JobName       string
	JobCmd        string
	WarningLines  int
}
```

TaskConfig defines the configuration for the task this program wraps.

#### func (TaskConfig) Process

```go
func (conf TaskConfig) Process(job baseworker.Job) (data []byte, err error)
```
Process runs the Gearman job by running the configured task.

## Testing

You can run the test cases by typing `make test` in the root of the repository

## Documentation

The documentation is automatically generated via [godocdown](https://github.com/robertkrimen/godocdown).

You can update it by typing `make docs` in the root of the repository

They're also viewable online at [![GoDoc](https://godoc.org/github.com/Clever/baseworker-go/workerwrapper?status.png)](https://godoc.org/github.com/Clever/baseworker-go/workwrapper).
