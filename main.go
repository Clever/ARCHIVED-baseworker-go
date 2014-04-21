package main

import (
	"github.com/Clever/baseworker-go/worker"
	"os"
)

func jobFunc(job gearman.Job) ([]byte, error) {
	return []byte{}, nil
}

func main() {
	worker := gearman.New("test", jobFunc)
	err := worker.Listen(os.Getenv("GEARMAN_HOST"), os.Getenv("GEARMAN_PORT"))
	if err != nil {
		panic(err)
	}
}
