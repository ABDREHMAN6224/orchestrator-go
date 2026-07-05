package main

import (
	"fmt"
	"time"

	"github.com/abdrehman6224/orchestrator/task"
	"github.com/abdrehman6224/orchestrator/worker"
	"github.com/google/uuid"
)

func main() {
	w := worker.NewWorker("test-1")
	t := task.Task{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}
	fmt.Println("starting task")
	w.AddTask(t)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
	t.ContainerID = result.ContainerId
	fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 10)
	fmt.Printf("stopping task %s\n", t.ID)
	t.State = task.Completed
	w.AddTask(t)
	result = w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
}
