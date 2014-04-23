# baseworker
--
    import "github.com/Clever/baseworker-go"

Package gearman provides a simple wrapper around a Gearman worker, based on
http://godoc.org/github.com/mikespook/gearman-go.


### Example

Here's an example program that just listens for "test" jobs and logs the data
that it receives:

    package main

    import(
    	"github.com/Clever/baseworker-go"
    	"log"
    )

    func jobFunc(job gearman.Job) ([]byte, error) {
    	log.Printf("Got job with data %s", job.Data())
    	return []byte{}, nil
    }

    func main() {
    	worker := gearman.New("test", jobFunc)
    	worker.Listen("localhost", "4730")
    }

## Usage

#### type Job

```go
type Job gearmanWorker.Job
```

Job is an alias for http://godoc.org/github.com/mikespook/gearman-go/worker#Job.

#### type JobFunc

```go
type JobFunc func(Job) ([]byte, error)
```

JobFunc is a function that takes in a Gearman job and does some work on it.

#### type Worker

```go
type Worker struct {
}
```

Worker represents a Gearman worker.

#### func  New

```go
func New(name string, fn JobFunc) *Worker
```
New creates a new gearman worker with the specified name and job function.

#### func (*Worker) Listen

```go
func (worker *Worker) Listen(host, port string) error
```
Listen starts listening for jobs on the specified host and port.

## Testing

INSERT TESTING INSTRUCTIONS

## Documentation

INSERT INSTRUCTIONS FOR GENERATING DOCUMENTATION - POSSIBLY MAKE IT HAPPEN AUTOMATICALLY?
