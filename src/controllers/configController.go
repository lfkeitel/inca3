package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

type Config struct {
	e *utils.Environment
}

func NewConfig(e *utils.Environment) *Config {
	return &Config{e: e}
}

func (c *Config) ApiGetConfig(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("config")

	ret := utils.NewAPIResponse("", nil)
	if name == "" { // Get all configs
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
	config, err := models.GetConfigByID(c.e, name)
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

func (c *Config) ApiGetDeviceConfigs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("id")
	configID := p.ByName("config")

	ret := utils.NewAPIResponse("", nil)
	if name == "" {
		ret.WriteResponse(w, http.StatusNotFound)
		return
	}

	// Check device exists
	device, err := models.GetDeviceByID(c.e, name)
	if err != nil {
		c.e.Log.WithField("error", err).Error("Couldn't get device")
		return
	}

	if device.ID == "" {
		ret.WriteResponse(w, http.StatusNotFound)
		return
	}

	if configID == "" { // Return all configs
		configs, err := models.GetConfigsForDevice(c.e, device.ID)
		if err != nil {
			c.e.Log.WithField("error", err).Error("Couldn't get configs")
			return
		}

		ret.Data = configs
		ret.WriteResponse(w, http.StatusOK)
		return
	}

	config, err := models.GetConfigByID(c.e, configID)
	if err != nil {
		c.e.Log.WithField("error", err).Error("Couldn't get configs")
		return
	}

	ret.Data = config
	ret.WriteResponse(w, http.StatusOK)
	return
}
