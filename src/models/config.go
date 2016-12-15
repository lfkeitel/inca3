package models

import (
	"compress/gzip"
	"encoding/json"
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
	DeviceID   int       `json:"deviceID"`
	Filename   string    `json:"-"`
	Created    time.Time `json:"created"`
	Compressed bool      `json:"compressed"`
	Size       int64     `json:"size"`
	Text       string    `json:"body"`
}

func newConfig(e *utils.Environment) *Config {
	return &Config{e: e}
}

func GetAllConfigs(e *utils.Environment) ([]*Config, error) {
	return doConfigQuery(e, "", nil)
}

func GetConfigsForDevice(e *utils.Environment, id int) ([]*Config, error) {
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
			&c.DeviceID,
			&created,
			&c.Filename,
			&c.Compressed,
		)
		if err != nil {
			continue
		}

		f := filepath.Join(e.Config.Configs.BaseDir, c.Filename)
		if !utils.FileExists(f) {
			continue
		}

		fileInto, err := os.Stat(f)
		if err != nil {
			continue
		}

		c.Size = fileInto.Size()
		c.Created = time.Unix(created, 0)
		results = append(results, c)
	}
	return results, nil
}

func (c *Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	return json.Marshal(&struct {
		Created int64 `json:"created"`
		*Alias
	}{
		Created: c.Created.Unix(),
		Alias:   (*Alias)(c),
	})
}

func (c *Config) LoadText() error {
	t, err := c.getText()
	if err != nil {
		return err
	}
	c.Text = string(t)
	return nil
}

func (c *Config) getText() ([]byte, error) {
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
