package models

import (
	"compress/gzip"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/lfkeitel/inca3/src/utils"
)

type Config struct {
	e          *utils.Environment
	ID         string    `json:"id"`
	Device     string    `json:"device"`
	Filename   string    `json:"filename"`
	Created    time.Time `json:"created"`
	Compressed bool      `json:"compressed"`
}

func newConfig(e *utils.Environment) *Config {
	return &Config{e: e}
}

func GetConfigsForDevice(e *utils.Environment, id string) ([]*Config, error) {
	return doConfigQuery(e, `WHERE "device" = ? ORDER BY "created" DESC`, id)
}

func GetConfigByID(e *utils.Environment, id string) (*Config, error) {
	configs, err := doConfigQuery(e, `WHERE "id" = ?`, id)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return newConfig(e), nil
	}
	return configs[0], nil
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

func (c *Config) GetText() ([]byte, error) {
	filename := filepath.Join(c.e.Config.Configs.BaseDir, c.Filename)

	if !utils.FileExists(filename) {
		return nil, errors.New("Config file not found")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if !c.Compressed {
		return ioutil.ReadAll(file)
	}

	// Uncompress file
	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
