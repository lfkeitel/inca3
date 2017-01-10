package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/utils"
)

var logControllerSingle *LogController

type LogController struct {
	e *utils.Environment
}

func newLogController(e *utils.Environment) *LogController {
	return &LogController{e: e}
}

func GetLogController(e *utils.Environment) *LogController {
	if logControllerSingle == nil {
		logControllerSingle = newLogController(e)
	}
	return logControllerSingle
}

func (l *LogController) ApiGetLogs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := utils.NewAPIResponse("", l.e.Log.GetUserLogs())
	response.WriteResponse(w, http.StatusOK)
}
