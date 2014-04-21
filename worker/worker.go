package gearman

import (
	"errors"
	"fmt"
	gearmanWorker "github.com/azylman/gearman-go/worker"
	"github.com/mikespook/golib/signal"
	"log"
	"net"
	"os"
)

// The same as a gearmanWorker.JobFunc, but takes in a gearman.Job instead of a gearmanWorker.Job.
type JobFunc func(Job) ([]byte, error)

// Alias for gearmanWorker.Job just so the consumer only needs one Gearman library.
type Job gearmanWorker.Job

// A Gearman worker.
type Worker struct {
	fn   gearmanWorker.JobFunc
	name string
}

// Starts listening
func (worker *Worker) Listen(host, port string) error {
	if host == "" || port == "" {
		return errors.New("must provide host and port")
	}
	w := gearmanWorker.New(1)
	defer w.Close()
	w.ErrorHandler = func(e error) {
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
	w.AddServer("tcp4", fmt.Sprintf("%s:%s", host, port))
	w.AddFunc(worker.name, worker.fn, gearmanWorker.Immediately)
	if err := w.Ready(); err != nil {
		log.Fatal(err)
	}
	go w.Work()
	sh := signal.NewHandler()
	sh.Bind(os.Interrupt, func() bool { return true })
	sh.Loop()
	return nil
}

// Creates a new gearman worker with the specified name and job function.
func New(name string, fn JobFunc) *Worker {
	// Turn a JobFunc into gearmanWorker.JobFunc
	jobFunc := func(job gearmanWorker.Job) ([]byte, error) {
		castedJob := Job(job)
		return fn(castedJob)
	}
	return &Worker{fn: jobFunc, name: name}
}
