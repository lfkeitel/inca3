package controllers

import (
	"net/http"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

var typeControllerSingle *TypeController

type TypeController struct {
	e *utils.Environment
}

func newTypeController(e *utils.Environment) *TypeController {
	return &TypeController{e: e}
}

func GetTypeController(e *utils.Environment) *TypeController {
	if typeControllerSingle == nil {
		typeControllerSingle = newTypeController(e)
	}
	return typeControllerSingle
}

func (t *TypeController) ApiGetTypes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resp := utils.NewAPIResponse("", nil)

	typeIDStr := p.ByName("id")

	if typeIDStr == "" {
		types, err := models.GetAllTypes(t.e)
		if err != nil {
			resp.Message = err.Error()
			resp.WriteResponse(w, http.StatusInternalServerError)
			return
		}
		resp.Data = types
		resp.WriteResponse(w, http.StatusOK)
		return
	}

	id, err := strconv.Atoi(typeIDStr)
	if err != nil {
		resp.Message = "Invalid type ID"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	dType, err := models.GetTypeByID(t.e, id)
	if err != nil {
		resp.Message = "Error loading type"
		resp.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	if dType.ID == 0 {
		resp.WriteResponse(w, http.StatusNotFound)
		return
	}

	resp.Data = dType
	resp.WriteResponse(w, http.StatusOK)
}

func (t *TypeController) ApiPutType(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := utils.NewAPIResponse("Not implemented", nil)
	resp.WriteResponse(w, http.StatusBadRequest)
}

func (t *TypeController) ApiPostType(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := utils.NewAPIResponse("Not implemented", nil)
	resp.WriteResponse(w, http.StatusBadRequest)
}

func (t *TypeController) ApiDeleteType(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := utils.NewAPIResponse("Not implemented", nil)
	resp.WriteResponse(w, http.StatusBadRequest)
}
