package main

import (
	"flag"
	"log"
	"os"

	baseworker "github.com/Clever/baseworker-go"
	workerwrapper "github.com/Clever/baseworker-go/workerwrapper"
)

func main() {
	JobName := flag.String("name", "", "Name of the Gearman job")
	JobCmd := flag.String("cmd", "", "The cmd to run")
	flag.Parse()
	if len(*JobName) == 0 {
		log.Printf("Error: name not defined")
		flag.PrintDefaults()
		os.Exit(2)
	}
	if len(*JobCmd) == 0 {
		log.Printf("Error: cmd not defined")
		flag.PrintDefaults()
		os.Exit(3)
	}

	config := workerwrapper.WorkerConfig{JobName: *JobName, JobCmd: *JobCmd, WarningLines: 5}
	worker := baseworker.NewWorker(*JobName, config.Process)
	defer worker.Close()
	log.Printf("Listing for job: " + *JobName)
	err := worker.Listen(os.Getenv("GEARMAN_HOST"), os.Getenv("GEARMAN_PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
