package models

import (
	"fmt"
	"strings"

	"github.com/lfkeitel/inca3/src/utils"
)

type Device struct {
	e       *utils.Environment
	ID      int          `json:"id"`
	Slug    string       `json:"slug"`
	Name    string       `json:"name"`
	Address string       `json:"address"`
	Profile *ConnProfile `json:"profile"`
	Configs []string     `json:"configs"`
}

func (d *Device) Print() {
	fmt.Printf("ID: %d\n", d.ID)
	fmt.Printf("Slug: %s\n", d.Slug)
	fmt.Printf("Name: %s\n", d.Name)
	fmt.Printf("Address: %s\n", d.Address)
	fmt.Printf("Type: %s\n", d.Profile.Name)
}

func NewDevice(e *utils.Environment) *Device {
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
		return NewDevice(e), nil
	}
	return devices[0], nil
}

func GetDeviceBySlug(e *utils.Environment, name string) (*Device, error) {
	devices, err := doDeviceQuery(e, `WHERE "slug" = ?`, name)
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return NewDevice(e), nil
	}
	return devices[0], nil
}

func GetDevicesByIDs(e *utils.Environment, ids []int) ([]*Device, error) {
	sql := "WHERE " + strings.Repeat(`"id" = ? OR `, len(ids))
	return doDeviceQuery(e, sql[:len(sql)-5], ids)
}

func doDeviceQuery(e *utils.Environment, where string, values ...interface{}) ([]*Device, error) {
	sql := `SELECT "id", "slug", "name", "address", "type" FROM "device" ` + where

	rows, err := e.DB.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Device
	for rows.Next() {
		var typeID int
		d := NewDevice(e)
		err := rows.Scan(
			&d.ID,
			&d.Slug,
			&d.Name,
			&d.Address,
			&typeID,
		)
		if err != nil {
			continue
		}
		d.loadConfigs()
		dType, err := GetTypeByID(e, typeID)
		if err != nil {
			return nil, err
		}
		d.Profile = dType
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
	if d.Profile.ID == 0 {
		if err := d.Profile.Save(); err != nil {
			return err
		}
	}
	d.Slug = utils.GenerateSlug(d.Name)

	if d.ID == 0 {
		return d.create()
	}
	return d.update()
}

func (d *Device) create() error {
	sql := `INSERT INTO "device" ("slug", "name", "address", "type") VALUES (?, ?, ?, ?)`

	result, err := d.e.DB.Exec(
		sql,
		d.Slug,
		d.Name,
		d.Address,
		d.Profile.ID,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	d.ID = int(id)
	return nil
}

func (d *Device) update() error {
	sql := `UPDATE "device" SET "slug" = ?, "name" = ?, "address" = ?, "type" = ? WHERE "id" = ?`

	_, err := d.e.DB.Exec(
		sql,
		d.Slug,
		d.Name,
		d.Address,
		d.Profile.ID,
		d.ID,
	)

	return err
}

func (d *Device) Delete() error {
	sql := `DELETE FROM "device" WHERE "id" = ?`
	_, err := d.e.DB.Exec(sql, d.ID)
	return err
}
