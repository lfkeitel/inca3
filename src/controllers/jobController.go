package controllers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/jobs"
	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

var jobControllerSingle *JobController

type JobController struct {
	e *utils.Environment
}

func newJobController(e *utils.Environment) *JobController {
	return &JobController{e: e}
}

func GetJobController(e *utils.Environment) *JobController {
	if jobControllerSingle == nil {
		jobControllerSingle = newJobController(e)
	}
	return jobControllerSingle
}

func (j *JobController) ApiJobStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idStr := p.ByName("id")
	resp := utils.NewAPIResponse("", nil)

	if idStr == "" {
		resp.Message = "No ID given"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp.Message = "Invalid ID"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	status, err := jobs.StatusJob(id)
	if err != nil {
		job, err := models.GetJobByID(j.e, id)
		if err != nil {
			j.e.Log.WithField("Err", err).Debug("Error getting job from database")
			resp.Message = "Error getting job status"
			resp.WriteResponse(w, http.StatusInternalServerError)
			return
		}

		if job == nil {
			resp.Message = "Job does not exist"
			resp.WriteResponse(w, http.StatusNotFound)
			return
		}

		status = &jobs.JobStatus{
			Started:   job.Start,
			Finished:  job.End,
			Total:     job.Total,
			Completed: job.Total,
		}
	}

	resp.Data = status
	resp.WriteResponse(w, http.StatusOK)
}

func (j *JobController) ApiStartJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	resp := utils.NewAPIResponse("", nil)
	var jobRequest *models.JobApiRequest
	err := decoder.Decode(&jobRequest)
	if err != nil {
		resp.Message = "Invalid JSON"
		j.e.Log.WithField("Err", err).Error("Invalid JSON")
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if jobRequest.Start != 0 {
		resp.Message = "Scheduled jobs are not supported"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	job := models.NewJob(j.e)
	job.Type = jobRequest.Type
	job.Devices = jobRequest.Devices
	job.Status = models.Pending

	jobid, err := jobs.StartJob(j.e, job)
	if err != nil {
		j.e.Log.WithField("Err", err).Error()
		resp.Message = "Job failed to start"
		resp.WriteResponse(w, http.StatusTeapot)
		return
	}

	resp.Data = struct {
		ID int `json:"id"`
	}{jobid}
	resp.WriteResponse(w, http.StatusAccepted)
}

func (j *JobController) ApiStopJob(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resp := utils.NewAPIResponse("Not implemented", nil)
	resp.WriteResponse(w, http.StatusNotImplemented)
}
