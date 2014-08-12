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
	JobExecutable := flag.String("exec", "", "The executable to run")
	flag.Parse()
	if len(*JobName) == 0 {
		log.Printf("Error: name not defined")
		flag.PrintDefaults()
		os.Exit(2)
	}
	if len(*JobExecutable) == 0 {
		log.Printf("Error: exec not defined")
		flag.PrintDefaults()
		os.Exit(3)
	}

	config := workerwrapper.WorkerConfig{JobName: *JobName, JobExecutable: *JobExecutable, WarningLines: 5}
	worker := baseworker.NewWorker(*JobName, config.Process)
	defer worker.Close()
	log.Printf("Listing for job: " + *JobName)
	err := worker.Listen(os.Getenv("GEARMAN_HOST"), os.Getenv("GEARMAN_PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
