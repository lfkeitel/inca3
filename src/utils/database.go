// This source file is part of the Inca project.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type databaseInit func(*DatabaseAccessor, *Config) error

var dbInits = make(map[string]databaseInit)

type DatabaseAccessor struct {
	*sql.DB
	Driver string
}

func NewDatabaseAccessor(config *Config) (*DatabaseAccessor, error) {
	d := &DatabaseAccessor{}
	var err error
	if err = os.MkdirAll(path.Dir(c.Database.Address), os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create directories: %v", err)
	}
	d.DB, err = sql.Open("sqlite3", c.Database.Address)
	if err != nil {
		return nil, err
	}

	err = d.DB.Ping()
	if err != nil {
		return nil, err
	}

	d.Driver = "sqlite"
	if _, err = d.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	rows, err := d.DB.Query(`SELECT name FROM sqlite_master WHERE type='table'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables := make(map[string]bool)

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables[tableName] = true
	}

	if _, ok := tables["device"]; !ok {
		if err := createDeviceTable(d); err != nil {
			return nil, err
		}
	}
	if _, ok := tables["type"]; !ok {
		if err := createTypeTable(d); err != nil {
			return nil, err
		}
	}
	if _, ok := tables["config"]; !ok {
		if err := createConfigTable(d); err != nil {
			return nil, err
		}
	}
	if _, ok := tables["log"]; !ok {
		if err := createLogTable(d); err != nil {
			return nil, err
		}
	}
	return d, nil
}

func createDeviceTable(d *databaseAccessor) error {
	return nil
}

func createTypeTable(d *databaseAccessor) error {
	return nil
}

func createConfigTable(d *databaseAccessor) error {
	return nil
}

func createLogTable(d *databaseAccessor) error {
	return nil
}
