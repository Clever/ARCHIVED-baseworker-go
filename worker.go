package baseworker

import (
	"errors"
	"fmt"
	gearmanWorker "github.com/azylman/gearman-go/worker"
	"log"
	"net"
	"os"
)

// JobFunc is a function that takes in a Gearman job and does some work on it.
type JobFunc func(Job) ([]byte, error)

// Job is an alias for http://godoc.org/github.com/mikespook/gearman-go/worker#Job.
type Job gearmanWorker.Job

// Worker represents a Gearman worker.
type Worker struct {
	fn   gearmanWorker.JobFunc
	name string
	w    *gearmanWorker.Worker
}

// Listen starts listening for jobs on the specified host and port.
func (worker *Worker) Listen(host, port string) error {
	if host == "" || port == "" {
		return errors.New("must provide host and port")
	}
	worker.w = gearmanWorker.New(1)
	worker.w.ErrorHandler = func(e error) {
		log.Println(e)
		if opErr, ok := e.(*net.OpError); ok {
			if !opErr.Temporary() {
				proc, err := os.FindProcess(os.Getpid())
				if err != nil {
					log.Println(err)
				}
				if err := proc.Signal(os.Interrupt); err != nil {
					log.Println(err)
				}
			}
		}
	}
	worker.w.AddServer("tcp4", fmt.Sprintf("%s:%s", host, port))
	worker.w.AddFunc(worker.name, worker.fn, gearmanWorker.Immediately)
	if err := worker.w.Ready(); err != nil {
		log.Fatal(err)
	}
	worker.w.Work()
	return nil
}

// Close closes the connection.
func (worker *Worker) Close() error {
	if worker.w == nil {
		return errors.New("no worker connection exists")
	}
	worker.w.Close()
	return nil
}

// New creates a new gearman worker with the specified name and job function.
func NewWorker(name string, fn JobFunc) *Worker {
	// Turn a JobFunc into gearmanWorker.JobFunc
	jobFunc := func(job gearmanWorker.Job) ([]byte, error) {
		castedJob := Job(job)
		return fn(castedJob)
	}
	return &Worker{fn: jobFunc, name: name}
}
