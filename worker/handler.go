package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/abdrehman6224/orchestrator/task"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}
	err := d.Decode(&te)
	if err != nil {
		msg := fmt.Sprintf("Error unmarshalling body: %v\n", err)
		log.Print(msg)
		w.WriteHeader(400)
		e := ErrResponse{
			HTTPStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
		return
	}
	a.Worker.AddTask(te.Task)
	log.Printf("added task %v\n", te.Task.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
}
func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}
func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	if taskId == "" {
		log.Println("No taskId is present in request")
		w.WriteHeader(400)
	}
	tId, _ := uuid.Parse(taskId)
	_, ok := a.Worker.Db[tId]
	if !ok {
		log.Printf("NO task found with ID= %v\n", tId)
		w.WriteHeader(404)
	}
	taskToStop := a.Worker.Db[tId]
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)
	log.Printf("Added task %v to stop container %v\n", taskToStop.ID,
		taskToStop.ContainerID)
	w.WriteHeader(204)
}
