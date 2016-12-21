package models

import "github.com/lfkeitel/inca3/src/utils"

type Type struct {
	e          *utils.Environment
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	Connection string `json:"connection"`
	Script     string `json:"script"`
	Args       string `json:"args"`
}

func newType(e *utils.Environment) *Type {
	return &Type{e: e}
}

func GetAllTypes(e *utils.Environment) ([]*Type, error) {
	return doTypeQuery(e, "", nil)
}

func GetTypeByID(e *utils.Environment, id int) (*Type, error) {
	types, err := doTypeQuery(e, `WHERE "id" = ?`, id)
	if err != nil {
		return nil, err
	}
	if len(types) == 0 {
		return newType(e), nil
	}
	return types[0], nil
}

func doTypeQuery(e *utils.Environment, where string, values ...interface{}) ([]*Type, error) {
	sql := `SELECT "id", "name", "brand", "connection", "script", "args" FROM "type" ` + where

	rows, err := e.DB.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Type
	for rows.Next() {
		t := newType(e)
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Brand,
			&t.Connection,
			&t.Script,
			&t.Args,
		)
		if err != nil {
			continue
		}
		results = append(results, t)
	}
	return results, nil
}

func (t *Type) Save() error {
	if t.ID == 0 {
		return t.create()
	}
	return t.update()
}

func (t *Type) create() error {
	sql := `INSERT INTO "type" ("name", "brand", "connection", "script", "args") VALUES (?,?,?,?,?)`

	result, err := t.e.DB.Exec(
		sql,
		t.Name,
		t.Brand,
		t.Connection,
		t.Script,
		t.Args,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	t.ID = int(id)
	return nil
}

func (t *Type) update() error {
	sql := `UPDATE "type" SET "name" = ?, "brand" = ?, "connection" = ?, "script" = ?, "args" = ? WHERE "id" = ?`

	_, err := t.e.DB.Exec(
		sql,
		t.Name,
		t.Brand,
		t.Connection,
		t.Script,
		t.Args,
		t.ID,
	)
	return err
}
