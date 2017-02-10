package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

var jm *jobManager

const maxWorkers = 5

type jobManager struct {
	e    *utils.Environment
	jobs map[int]*job
	//workers []*worker
}

type job struct {
	*models.Job
	worker   *worker
	finished int
}

type JobStatus struct {
	Started   time.Time `json:"started"`
	Finished  time.Time `json:"finished"`
	Total     int       `json:"total"`
	Completed int       `json:"completed"`
}

func (j *JobStatus) MarshalJSON() ([]byte, error) {
	type Alias JobStatus
	return json.Marshal(&struct {
		Started  int64 `json:"started"`
		Finished int64 `json:"finished"`
		*Alias
	}{
		Started:  j.Started.Unix(),
		Finished: j.Finished.Unix(),
		Alias:    (*Alias)(j),
	})
}

func getJobManager(e *utils.Environment) *jobManager {
	if jm == nil {
		jm = &jobManager{
			e:    e,
			jobs: make(map[int]*job),
		}
	}
	return jm
}

func StartJob(e *utils.Environment, j *models.Job) (int, error) {
	jm := getJobManager(e)
	if j.ID == 0 {
		if err := j.Save(); err != nil {
			return 0, err
		}
	}

	newJob := &job{Job: j}

	jm.jobs[newJob.ID] = newJob

	e.Log.Debug("Creating new worker")
	w := newWorker(e)
	w.job = newJob
	newJob.worker = w
	w.running = true

	go func() {
		e.Log.Debug("Waiting for worker to end")
		<-w.done
		e.Log.Debug("Worker done, deleting job")
		delete(jm.jobs, newJob.ID)
		w.running = false
	}()

	e.Log.Debug("Starting worker")
	return newJob.ID, w.start()
}

func StopJob(id int) error {
	return nil
}

func StatusJob(id int) (*JobStatus, error) {
	j, ok := getJobManager(nil).jobs[id]
	if !ok {
		return nil, fmt.Errorf("No job with ID %d", id)
	}

	return &JobStatus{
		Started:   j.Start,
		Total:     j.Total,
		Completed: j.finished,
		Finished:  j.End,
	}, nil
}
