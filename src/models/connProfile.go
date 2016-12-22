package models

import "github.com/lfkeitel/inca3/src/utils"

type ConnProfile struct {
	e        *utils.Environment
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Username string `json:"username"`
	Password string `json:"password"`
	EnablePW string `json:"enablepw"`
	Script   string `json:"script"`
}

func NewType(e *utils.Environment) *ConnProfile {
	return &ConnProfile{e: e}
}

func GetAllTypes(e *utils.Environment) ([]*ConnProfile, error) {
	return doTypeQuery(e, "", nil)
}

func GetTypeBySlug(e *utils.Environment, name string) (*ConnProfile, error) {
	types, err := doTypeQuery(e, `WHERE "slug" = ?`, name)
	if err != nil || len(types) == 0 {
		return nil, err
	}
	return types[0], nil
}

func GetTypeByID(e *utils.Environment, id int) (*ConnProfile, error) {
	types, err := doTypeQuery(e, `WHERE "id" = ?`, id)
	if err != nil {
		return nil, err
	}
	if len(types) == 0 {
		return nil, nil
	}
	return types[0], nil
}

func doTypeQuery(e *utils.Environment, where string, values ...interface{}) ([]*ConnProfile, error) {
	sql := `SELECT "id", "name", "slug", "username", "password", "enable", "script" FROM "type" ` + where

	rows, err := e.DB.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*ConnProfile
	for rows.Next() {
		t := NewType(e)
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Slug,
			&t.Username,
			&t.Password,
			&t.EnablePW,
			&t.Script,
		)
		if err != nil {
			continue
		}
		results = append(results, t)
	}
	return results, nil
}

func (t *ConnProfile) Save() error {
	t.Slug = utils.GenerateSlug(t.Name)

	if t.ID == 0 {
		return t.create()
	}
	return t.update()
}

func (t *ConnProfile) create() error {
	sql := `INSERT INTO "type" ("name", "slug", "username", "password", "enable", "script") VALUES (?,?,?,?,?,?)`

	result, err := t.e.DB.Exec(
		sql,
		t.Name,
		t.Slug,
		&t.Username,
		&t.Password,
		&t.EnablePW,
		t.Script,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	t.ID = int(id)
	return nil
}

func (t *ConnProfile) update() error {
	sql := `UPDATE "type" SET "name" = ?, "slug" = ?, "username" = ?, "password" = ?, "enable" = ?, "script" = ? WHERE "id" = ?`

	_, err := t.e.DB.Exec(
		sql,
		t.Name,
		t.Slug,
		&t.Username,
		&t.Password,
		&t.EnablePW,
		t.Script,
		t.ID,
	)
	return err
}

func (t *ConnProfile) Delete() error {
	sql := `DELETE FROM "type" WHERE "id" = ?`
	_, err := t.e.DB.Exec(sql, t.ID)
	return err
}
