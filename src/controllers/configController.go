package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

var configControllerSingle *ConfigController

type ConfigController struct {
	e *utils.Environment
}

func newConfigController(e *utils.Environment) *ConfigController {
	return &ConfigController{e: e}
}

func GetConfigController(e *utils.Environment) *ConfigController {
	if configControllerSingle == nil {
		configControllerSingle = newConfigController(e)
	}
	return configControllerSingle
}

func (c *ConfigController) ApiGetConfig(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	configSlug := p.ByName("config")

	ret := utils.NewAPIResponse("", nil)
	if configSlug == "" { // Get all configs
		configs, err := models.GetAllConfigs(c.e)
		if err != nil {
			ret.Message = "Error getting configs"
			ret.WriteResponse(w, http.StatusInternalServerError)
			return
		}

		for _, config := range configs {
			if err := config.LoadText(); err != nil {
				ret.Message = "Error loading config from file"
				ret.WriteResponse(w, http.StatusInternalServerError)
				return
			}
		}

		ret.Data = configs
		ret.WriteResponse(w, http.StatusOK)
		return
	}

	// Get a single config
	config, err := models.GetConfigBySlug(c.e, configSlug)
	if err != nil {
		ret.Message = "Error getting config"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	if err := config.LoadText(); err != nil {
		ret.Message = "Error loading config from file"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	ret.Data = config
	ret.WriteResponse(w, http.StatusOK)
	return
}

func (c *ConfigController) ApiGetDeviceConfigs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("slug")
	configSlug := p.ByName("config")

	ret := utils.NewAPIResponse("", nil)
	if name == "" {
		ret.WriteResponse(w, http.StatusNotFound)
		return
	}

	// Check device exists
	device, err := models.GetDeviceBySlug(c.e, name)
	if err != nil {
		c.e.Log.WithField("error", err).Error("Couldn't get device")
		return
	}

	if device.ID == 0 {
		ret.WriteResponse(w, http.StatusNotFound)
		return
	}

	if configSlug == "" { // Return all configs
		configs, err := models.GetConfigsForDevice(c.e, device.ID)
		if err != nil {
			c.e.Log.WithField("error", err).Error("Couldn't get configs")
			return
		}

		ret.Data = configs
		ret.WriteResponse(w, http.StatusOK)
		return
	}

	config, err := models.GetConfigBySlug(c.e, configSlug)
	if err != nil {
		c.e.Log.WithField("error", err).Error("Couldn't get configs")
		return
	}

	ret.Data = config
	ret.WriteResponse(w, http.StatusOK)
	return
}
