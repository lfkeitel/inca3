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
	ID         int       `json:"id"`
	Slug       string    `json:"slug"`
	DeviceID   int       `json:"deviceID"`
	Filename   string    `json:"-"`
	Created    time.Time `json:"created"`
	Compressed bool      `json:"compressed"`
	Size       int64     `json:"size"`
	Text       string    `json:"body"`
	Failed     bool      `json:"failed"`
}

func NewConfig(e *utils.Environment) *Config {
	return &Config{e: e}
}

func GetAllConfigs(e *utils.Environment) ([]*Config, error) {
	return doConfigQuery(e, "", nil)
}

func GetConfigsForDevice(e *utils.Environment, id int) ([]*Config, error) {
	return doConfigQuery(e, `WHERE "device" = ? ORDER BY "created" DESC`, id)
}

func GetConfigBySlug(e *utils.Environment, slug string) (*Config, error) {
	configs, err := doConfigQuery(e, `WHERE "slug" = ?`, slug)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return NewConfig(e), nil
	}
	return configs[0], nil
}

func GetConfigByID(e *utils.Environment, id string) (*Config, error) {
	configs, err := doConfigQuery(e, `WHERE "id" = ?`, id)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return NewConfig(e), nil
	}
	return configs[0], nil
}

func doConfigQuery(e *utils.Environment, where string, values ...interface{}) ([]*Config, error) {
	sql := `SELECT "id", "slug", "device", "created", "filename", "compressed", "failed" FROM "config" ` + where

	rows, err := e.DB.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cleanUpConfigs []*Config
	var results []*Config
	for rows.Next() {
		c := NewConfig(e)
		var created int64
		err := rows.Scan(
			&c.ID,
			&c.Slug,
			&c.DeviceID,
			&created,
			&c.Filename,
			&c.Compressed,
			&c.Failed,
		)
		if err != nil {
			continue
		}

		f := filepath.Join(e.Config.DirPaths.BaseDir, c.Filename)
		if !utils.FileExists(f) {
			cleanUpConfigs = append(cleanUpConfigs, c)
			continue
		}

		fileInto, err := os.Stat(f)
		if err != nil {
			e.Log.WithField("File", f).Warning("Config file not readable")
			continue
		}

		c.Size = fileInto.Size()
		c.Created = time.Unix(created, 0)
		results = append(results, c)
	}

	rows.Close()

	for _, c := range cleanUpConfigs {
		f := filepath.Join(e.Config.DirPaths.BaseDir, c.Filename)
		e.Log.WithField("File", f).Warning("Config file doesn't exist, removing from database")
		if err := c.Delete(); err != nil {
			e.Log.WithField("Err", err).Error("Failed to delete config from database")
		}
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
	filename := filepath.Join(c.e.Config.DirPaths.BaseDir, c.Filename)

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

func (c *Config) Save() error {
	if c.ID == 0 {
		return c.create()
	}
	return c.update()
}

func (c *Config) create() error {
	sql := `INSERT INTO "config" ("slug", "device", "created", "filename", "compressed", "failed") VALUES (?,?,?,?,?,?)`

	result, err := c.e.DB.Exec(
		sql,
		c.Slug,
		c.DeviceID,
		c.Created.Unix(),
		c.Filename,
		c.Compressed,
		c.Failed,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	c.ID = int(id)
	return nil
}

func (c *Config) update() error {
	sql := `UPDATE "config" SET "slug" = ?, "device" = ?, "created" = ?, "filename" = ?, "compressed" = ?, "failed" = ? WHERE "id" = ?`

	_, err := c.e.DB.Exec(
		sql,
		c.Slug,
		c.DeviceID,
		c.Created.Unix(),
		c.Filename,
		c.Compressed,
		c.Failed,
		c.ID,
	)
	return err
}

func (c *Config) Delete() error {
	// Delete the database record
	sql := `DELETE FROM "config" WHERE "id" = ?`
	if _, err := c.e.DB.Exec(sql, c.ID); err != nil {
		return err
	}

	// Delete the file
	f := filepath.Join(c.e.Config.DirPaths.BaseDir, c.Filename)
	if utils.FileExists(f) {
		return os.Remove(f)
	}
	return nil
}
