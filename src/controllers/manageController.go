package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/utils"
)

type Manager struct {
	e *utils.Environment
}

func NewManager(e *utils.Environment) *Manager {
	return &Manager{e: e}
}

func (m *Manager) Manage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// If object == devices, show the device management page
	// If object == types, show the type management page
}
