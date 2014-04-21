/*
Package gearman provides a simple wrapper around a Gearman worker, based on
http://godoc.org/github.com/mikespook/gearman-go.

Example

Here's an example program that just listens for "test" jobs and logs the data that it receives:

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
*/
package gearman
