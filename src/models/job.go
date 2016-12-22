package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/lfkeitel/inca3/src/utils"
)

type JobStatus string

const (
	Pending  JobStatus = "pending"
	Starting JobStatus = "starting"
	Running  JobStatus = "running"
	Stopping JobStatus = "stopping"
	Stopped  JobStatus = "stopped"
	Finished JobStatus = "finished"
)

type Job struct {
	e        *utils.Environment
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Status   JobStatus `json:"status"`
	Devices  []int     `json:"devices"`
	Total    int       `json:"total"`
	Finished int       `json:"finished"`
	Error    string    `json:"error"`
	Start    time.Time `json:"startTime"`
	End      time.Time `json:"endTime"`
}

func NewJob(e *utils.Environment) *Job {
	return &Job{e: e}
}

func GetAllJobs(e *utils.Environment) ([]*Job, error) {
	return doJobQuery(e, "", nil)
}

func GetJobsForDevice(e *utils.Environment, slug string) ([]*Job, error) {
	slug = "%;" + slug + ";%"
	return doJobQuery(e, `WHERE "devices" LIKE ? ORDER BY "start" DESC`, slug)
}

func GetJobByID(e *utils.Environment, id int) (*Job, error) {
	job, err := doJobQuery(e, `WHERE "id" = ?`, id)
	if err != nil {
		return nil, err
	}

	if len(job) == 0 {
		return nil, nil
	}
	return job[0], nil
}

func doJobQuery(e *utils.Environment, where string, values ...interface{}) ([]*Job, error) {
	sql := `SELECT "id", "name", "status", "type", "devices", "total", "error", "start", "end" FROM "job" ` + where

	rows, err := e.DB.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Job
	for rows.Next() {
		j := NewJob(e)
		var start int64
		var end int64
		var devices string
		var status string
		err := rows.Scan(
			&j.ID,
			&j.Name,
			&status,
			&j.Type,
			&devices,
			&j.Total,
			&j.Error,
			&start,
			&end,
		)
		if err != nil {
			e.Log.WithField("Err", err).Debug("Error scanning job struct")
			continue
		}

		j.Status = JobStatus(status)
		j.Start = time.Unix(start, 0)
		j.End = time.Unix(end, 0)
		if devices != "" {
			if err := parseDeviceInts(j, devices); err != nil {
				return nil, err
			}
		}
		results = append(results, j)
	}
	return results, nil
}

func parseDeviceInts(j *Job, devs string) error {
	devIDsStr := strings.Split(devs, ";")
	for _, id := range devIDsStr {
		idint, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		j.Devices = append(j.Devices, idint)
	}
	return nil
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

func (j *Job) Save() error {
	if j.ID == 0 {
		return j.create()
	}
	return j.update()
}

func (j *Job) create() error {
	sql := `INSERT INTO "job" ("name", "status", "type", "devices", "total", "error", "start", "end") VALUES (?,?,?,?,?,?,?,?)`

	result, err := j.e.DB.Exec(
		sql,
		j.Name,
		string(j.Status),
		j.Type,
		strings.Join(utils.IntSliceToString(j.Devices), ";"),
		j.Total,
		j.Error,
		j.Start.Unix(),
		j.End.Unix(),
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	j.ID = int(id)
	return nil
}

func (j *Job) update() error {
	sql := `UPDATE "job" SET "name" = ?, "status" = ?, "type" = ?, "devices" = ?, "total" = ?, "error" = ?, "start" = ?, "end" = ? WHERE "id" = ?`

	_, err := j.e.DB.Exec(
		sql,
		j.Name,
		string(j.Status),
		j.Type,
		strings.Join(utils.IntSliceToString(j.Devices), ";"),
		j.Total,
		j.Error,
		j.Start.Unix(),
		j.End.Unix(),
		j.ID,
	)
	return err
}
