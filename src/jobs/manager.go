package jobs

import "github.com/lfkeitel/inca3/src/models"
import "github.com/lfkeitel/inca3/src/utils"

var jm *jobManager

type jobManager struct {
	e       *utils.Environment
	jobs    map[int]*models.Job
	workers []*worker
}

func getJobManager(e *utils.Environment) *jobManager {
	if jm == nil {
		jm = &jobManager{
			e:    e,
			jobs: make(map[int]*models.Job),
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

	jm.jobs[j.ID] = j

	w := newWorker(e)
	w.job = j

	err := w.start()
	return j.ID, err
}

func StopJob(id int) error {
	return nil
}

func StatusJob(id int) (*models.JobStatus, error) {
	return nil, nil
}
