package models

import (
	"strings"

	"github.com/lfkeitel/inca3/src/utils"
)

type Device struct {
	e          *utils.Environment
	ID         int      `json:"id"`
	Slug       string   `json:"slug"`
	Name       string   `json:"name"`
	Address    string   `json:"address"`
	Brand      string   `json:"brand"`
	Connection string   `json:"connection"`
	Configs    []string `json:"configs"`
}

func newDevice(e *utils.Environment) *Device {
	return &Device{e: e}
}

func GetAllDevices(e *utils.Environment) ([]*Device, error) {
	return doDeviceQuery(e, "", nil)
}

func GetDeviceByID(e *utils.Environment, id int) (*Device, error) {
	devices, err := doDeviceQuery(e, `WHERE "id" = ?`, id)
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return newDevice(e), nil
	}
	return devices[0], nil
}

func GetDeviceBySlug(e *utils.Environment, name string) (*Device, error) {
	devices, err := doDeviceQuery(e, `WHERE "slug" = ?`, name)
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return newDevice(e), nil
	}
	return devices[0], nil
}

func doDeviceQuery(e *utils.Environment, where string, values ...interface{}) ([]*Device, error) {
	sql := `SELECT "id", "slug", "name", "address", "brand", "connection" FROM "device" ` + where

	rows, err := e.DB.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Device
	for rows.Next() {
		d := newDevice(e)
		err := rows.Scan(
			&d.ID,
			&d.Slug,
			&d.Name,
			&d.Address,
			&d.Brand,
			&d.Connection,
		)
		if err != nil {
			continue
		}
		d.loadConfigs()
		results = append(results, d)
	}
	return results, nil
}

func (d *Device) SetEnv(e *utils.Environment) {
	d.e = e
}

func (d *Device) loadConfigs() error {
	sql := `SELECT "id" FROM "config" WHERE "device" = ?`

	rows, err := d.e.DB.Query(sql, d.ID)
	if err != nil {
		return err
	}

	d.Configs = make([]string, 0)
	for rows.Next() {
		var c string
		err := rows.Scan(&c)
		if err != nil {
			continue
		}
		d.Configs = append(d.Configs, c)
	}
	return nil
}

func (d *Device) Save() error {
	d.Slug = d.generateSlug(d.Name)

	if d.ID == 0 {
		return d.create()
	}
	return d.update()
}

func (d *Device) generateSlug(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.Title(raw)
	return strings.Replace(raw, " ", "", -1)
}

func (d *Device) create() error {
	sql := `INSERT INTO "device" ("slug", "name", "address", "brand", "connection") VALUES (?, ?, ?, ?, ?)`

	result, err := d.e.DB.Exec(
		sql,
		d.Slug,
		d.Name,
		d.Address,
		d.Brand,
		d.Connection,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	d.ID = int(id)
	return nil
}

func (d *Device) update() error {
	sql := `UPDATE "device" SET "slug" = ?, "name" = ?, "address" = ?, "brand" = ?, "connection" = ? WHERE "id" = ?`

	_, err := d.e.DB.Exec(
		sql,
		d.Slug,
		d.Name,
		d.Address,
		d.Brand,
		d.Connection,
		d.ID,
	)

	return err
}
