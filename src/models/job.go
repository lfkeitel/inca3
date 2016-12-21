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

func doJobQuery(e *utils.Environment, where string, values ...interface{}) ([]*Job, error) {
	sql := `SELECT "id", "name", "status", "type", "devices", "error", "start", "end" FROM "job" ` + where

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
		err := rows.Scan(
			&j.ID,
			&j.Name,
			&j.Status,
			&j.Type,
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
		devIDsStr := strings.Split(devices, ";")
		for _, id := range devIDsStr {
			idint, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			j.Devices = append(j.Devices, idint)
		}
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

func (j *Job) Save() error {
	if j.ID == 0 {
		return j.create()
	}
	return j.update()
}

func (j *Job) create() error {
	sql := `INSERT INTO "job" ("name", "status", "type", "devices", "error", "start", "end") VALUES (?,?,?,?,?,?,?)`

	result, err := j.e.DB.Exec(
		sql,
		j.Name,
		string(j.Status),
		j.Type,
		strings.Join(utils.IntSliceToString(j.Devices), ";"),
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
	sql := `UPDATE "job" SET "name" = ?, "status" = ?, "type" = ?, "devices" = ?, "error" = ?, "start" = ?, "end" = ? WHERE "id" = ?`

	_, err := j.e.DB.Exec(
		sql,
		j.Name,
		string(j.Status),
		j.Type,
		strings.Join(utils.IntSliceToString(j.Devices), ";"),
		j.Error,
		j.Start.Unix(),
		j.End.Unix(),
		j.ID,
	)
	return err
}
