package models

import "github.com/lfkeitel/inca3/src/utils"

type Device struct {
	e                                    *utils.Environment
	ID, Name, Address, Brand, Connection string
}

func newDevice(e *utils.Environment) *Device {
	return &Device{e: e}
}

func GetAllDevices(e *utils.Environment) ([]*Device, error) {
	return doDeviceQuery(e, "", nil)
}

func GetDeviceByID(e *utils.Environment, name string) (*Device, error) {
	devices, err := doDeviceQuery(e, `WHERE "id" = ?`, name)
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return newDevice(e), nil
	}
	return devices[0], nil
}

func doDeviceQuery(e *utils.Environment, where string, values ...interface{}) ([]*Device, error) {
	sql := `SELECT "id", "name", "address", "brand", "connection" FROM "device" ` + where

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
			&d.Name,
			&d.Address,
			&d.Brand,
			&d.Connection,
		)
		if err != nil {
			continue
		}
		results = append(results, d)
	}
	return results, nil
}
