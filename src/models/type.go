package models

import "github.com/lfkeitel/inca3/src/utils"

type Type struct {
	e                                     *utils.Environment
	ID                                    int
	Name, Brand, Connection, Script, Args string
}
