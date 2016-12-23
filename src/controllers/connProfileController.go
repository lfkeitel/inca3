package controllers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"io/ioutil"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

var connProfileControllerSingle *ConnProfileController

type ConnProfileController struct {
	e *utils.Environment
}

func newConnProfileController(e *utils.Environment) *ConnProfileController {
	return &ConnProfileController{e: e}
}

func GetConnProfileController(e *utils.Environment) *ConnProfileController {
	if connProfileControllerSingle == nil {
		connProfileControllerSingle = newConnProfileController(e)
	}
	return connProfileControllerSingle
}

func (t *ConnProfileController) ShowTypeList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	types, err := models.GetAllTypes(t.e)
	if err != nil {
		t.e.Log.WithField("Err", err).Error("Failed getting device types")
		return
	}

	data := map[string]interface{}{
		"types": types,
	}
	t.e.View.NewView("type-list", r).Render(w, data)
}

func (t *ConnProfileController) ShowType(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	slug := p.ByName("slug")

	if slug == "" {
		http.Redirect(w, r, t.e.Config.Core.SiteBasePath+"/types", http.StatusTemporaryRedirect)
		return
	}

	dType, err := models.GetTypeBySlug(t.e, slug)
	if err != nil {
		t.e.Log.WithField("Err", err).Error()
		http.Redirect(w, r, t.e.Config.Core.SiteBasePath+"/types", http.StatusTemporaryRedirect)
		return
	}

	data := map[string]interface{}{
		"type": dType,
	}
	t.e.View.NewView("type", r).Render(w, data)
}

func (t *ConnProfileController) ApiGetTypes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resp := utils.NewAPIResponse("", nil)
	slug := p.ByName("slug")

	if slug == "" {
		types, err := models.GetAllTypes(t.e)
		if err != nil {
			resp.Message = err.Error()
			resp.WriteResponse(w, http.StatusInternalServerError)
			return
		}
		resp.Data = types
		resp.WriteResponse(w, http.StatusOK)
		return
	} else if slug == "_scripts" {
		var scripts []string
		files, err := ioutil.ReadDir(t.e.Config.DirPaths.ScriptDir)
		if err != nil {
			resp.Message = "Failed to get script list"
			t.e.Log.WithField("Err", err).Error()
			resp.WriteResponse(w, http.StatusInternalServerError)
			return
		}

		// Find all executable files in the directory
		for _, file := range files {
			if file.Mode().Perm()&0111 != 0 {
				scripts = append(scripts, file.Name())
			}
		}

		resp.Data = scripts
		resp.WriteResponse(w, http.StatusOK)
		return
	}

	dType, err := models.GetTypeBySlug(t.e, slug)
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

func (t *ConnProfileController) ApiPutType(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	slug := p.ByName("slug")
	resp := utils.NewAPIResponse("", nil)

	if slug == "" {
		resp.Message = "No slug given"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var apiType *models.ConnProfile
	if err := decoder.Decode(&apiType); err != nil {
		resp.Message = "Invalid JSON"
		t.e.Log.WithField("Err", err).Error("Invalid JSON")
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	oldType, err := models.GetTypeBySlug(t.e, slug)
	if err != nil {
		t.e.Log.WithField("Err", err).Error()
		resp.Message = "Error getting type"
		resp.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	// Type doesn't exist
	if oldType == nil {
		resp.WriteResponse(w, http.StatusNotFound)
		return
	}

	if oldType.Name != "" {
		oldType.Name = apiType.Name
	}

	if oldType.Username != "" {
		oldType.Username = apiType.Username
	}

	if oldType.Password != "" {
		oldType.Password = apiType.Password
	}

	if oldType.EnablePW != "" {
		oldType.EnablePW = apiType.EnablePW
	}

	if oldType.Script != "" {
		if exists := utils.FileExists(filepath.Join(t.e.Config.DirPaths.ScriptDir, apiType.Script)); !exists {
			resp.Message = "Script files does not exist"
			resp.WriteResponse(w, http.StatusBadRequest)
			return
		}
		oldType.Script = apiType.Script
	}

	if err := oldType.Save(); err != nil {
		resp.Message = "Failed to save type"
		t.e.Log.WithField("Err", err).Error("Failed to save type")
		resp.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	resp.Message = "Type saved successfully"
	resp.Data = oldType
	resp.WriteResponse(w, http.StatusOK)
}

func (t *ConnProfileController) ApiPostType(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	resp := utils.NewAPIResponse("", nil)
	var apiType *models.ConnProfile
	if err := decoder.Decode(&apiType); err != nil {
		resp.Message = "Invalid JSON"
		t.e.Log.WithField("Err", err).Error("Invalid JSON")
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	postedType := models.NewType(t.e)

	postedType.Name = apiType.Name
	postedType.Username = apiType.Username
	postedType.Password = apiType.Password
	postedType.EnablePW = apiType.EnablePW
	postedType.Script = apiType.Script

	if err := postedType.Save(); err != nil {
		resp.Message = "Failed to save type"
		t.e.Log.WithField("Err", err).Error("Failed to save type")
		resp.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	resp.Message = "Type saved successfully"
	resp.Data = postedType
	resp.WriteResponse(w, http.StatusOK)
}

func (t *ConnProfileController) ApiDeleteType(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	slug := p.ByName("slug")

	ret := utils.NewAPIResponse("", nil)
	if slug == "" {
		ret.Message = "No slug given"
		ret.WriteResponse(w, http.StatusBadRequest)
		return
	}

	dType, err := models.GetTypeBySlug(t.e, slug)
	if err != nil {
		t.e.Log.WithField("Err", err).Error()
		ret.Message = "Error getting type"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	// Type doesn't exist
	if dType == nil {
		ret.WriteResponse(w, http.StatusNoContent)
		return
	}

	if err := dType.Delete(); err != nil {
		t.e.Log.WithField("Err", err).Error()
		ret.Message = "Error deleting type"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	ret.Message = "Type deleted"
	ret.WriteResponse(w, http.StatusOK)
}
