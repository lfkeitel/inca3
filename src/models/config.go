package models

import (
	"time"

	"github.com/lfkeitel/inca3/src/utils"
)

type Config struct {
	e                    *utils.Environment
	ID, Device, Filename string
	Created              time.Time
	Compressed           bool
}

func newConfig(e *utils.Environment) *Config {
	return &Config{e: e}
}

func GetConfigsForDevice(e *utils.Environment, id string) ([]*Config, error) {
	return doConfigQuery(e, `WHERE "device" = ? ORDER BY "created" DESC`, id)
}

func doConfigQuery(e *utils.Environment, where string, values ...interface{}) ([]*Config, error) {
	sql := `SELECT "id", "device", "created", "filename", "compressed" FROM "config" ` + where

	rows, err := e.DB.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Config
	for rows.Next() {
		c := newConfig(e)
		var created int64
		err := rows.Scan(
			&c.ID,
			&c.Device,
			&created,
			&c.Filename,
			&c.Compressed,
		)
		if err != nil {
			continue
		}
		c.Created = time.Unix(created, 0)
		results = append(results, c)
	}
	return results, nil
}
