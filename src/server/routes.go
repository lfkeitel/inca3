package server

import (
	"net/http"

	"github.com/lfkeitel/inca3/src/utils"
)

func LoadRoutes(e *utils.Environment) http.Handler {
	return http.FileServer(http.Dir("public"))
}
