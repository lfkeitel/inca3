package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/lfkeitel/inca3/src/utils"
)

type JobStatus string

const (
	Started  JobStatus = "started"
	Running  JobStatus = "running"
	Stopped  JobStatus = "stopped"
	Finished JobStatus = "finished"
)

type Job struct {
	e        *utils.Environment
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Status   JobStatus `json:"status"`
	Devices  []string  `json:"devices"`
	Finished int       `json:"finished"`
	Error    string    `json:"error"`
	Start    time.Time `json:"startTime"`
	End      time.Time `json:"endTime"`
}

func newJob(e *utils.Environment) *Job {
	return &Job{e: e}
}

func GetAllJobs(e *utils.Environment) ([]*Job, error) {
	return doJobQuery(e, "", nil)
}

func GetJobsForDevice(e *utils.Environment, slug string) ([]*Job, error) {
	slug = "%;" + slug + ";%"
	return doJobQuery(e, `WHERE "devices" LIKE ? ORDER BY "start" DESC`, slug)
}

func doJobQuery(e *utils.Environment, where string, values ...interface{}) ([]*Job, error) {
	sql := `SELECT "id", "name", "status", "devices", "error", "start", "end" FROM "job" ` + where

	rows, err := e.DB.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Job
	for rows.Next() {
		j := newJob(e)
		var start int64
		var end int64
		var devices string
		err := rows.Scan(
			&j.ID,
			&j.Name,
			&j.Status,
			&devices,
			&j.Error,
			&start,
			&end,
		)
		if err != nil {
			continue
		}

		j.Start = time.Unix(start, 0)
		j.End = time.Unix(end, 0)
		j.Devices = strings.Split(devices, ";")
		results = append(results, j)
	}
	return results, nil
}

func (j *Job) MarshalJSON() ([]byte, error) {
	type Alias Job
	return json.Marshal(&struct {
		Start int64 `json:"startTime"`
		End   int64 `json:"endTime"`
		*Alias
	}{
		Start: j.Start.Unix(),
		End:   j.End.Unix(),
		Alias: (*Alias)(j),
	})
}
