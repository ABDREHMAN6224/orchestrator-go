package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/abdrehman6224/orchestrator/worker"
)

func main() {
	host := os.Getenv("ORCH_HOST")
	port, _ := strconv.Atoi(os.Getenv("ORCH_PORT"))
	fmt.Println("starting orchestrator")
	w := worker.NewWorker("test-1")
	api := worker.Api{
		Address: host,
		Port:    port,
		Worker:  w,
	}
	go runTasks(w)
	go w.CollectStats()
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}
