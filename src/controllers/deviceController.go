package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

type Device struct {
	e *utils.Environment
}

func NewDevice(e *utils.Environment) *Device {
	return &Device{e: e}
}

func (d *Device) ShowDevice(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if p.ByName("slug") != "" {
		d.showDeviceConfigList(w, r)
		return
	}

	d.e.View.NewView("device-list", r).Render(w, nil)
	return
}

func (d *Device) showDeviceConfigList(w http.ResponseWriter, r *http.Request) {
	d.e.View.NewView("device", r).Render(w, nil)
}

func (d *Device) ShowConfig(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("slug")
	configSlug := p.ByName("config")

	if configSlug == "" {
		d.showDeviceConfigList(w, r)
		return
	}

	device, err := models.GetDeviceBySlug(d.e, name)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get device")
		return
	}

	if device.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	config, err := models.GetConfigBySlug(d.e, configSlug)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get config")
		return
	}

	if config.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := config.LoadText(); err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get config text")
		return
	}

	data := map[string]interface{}{
		"device": device,
		"config": config,
	}
	d.e.View.NewView("config", r).Render(w, data)
}

func (d *Device) ApiGetDevices(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("slug")

	ret := utils.NewAPIResponse("", nil)
	if name == "" {
		devices, err := models.GetAllDevices(d.e)
		if err != nil {
			ret.Message = "Error getting devices"
			ret.WriteResponse(w, http.StatusInternalServerError)
			return
		}

		ret.Data = devices
		ret.WriteResponse(w, http.StatusOK)
		return
	}

	device, err := models.GetDeviceBySlug(d.e, name)
	if err != nil {
		ret.Message = "Error getting devices"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	ret.Data = device
	ret.WriteResponse(w, http.StatusOK)
	return
}

func (d *Device) ApiSaveDevice(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	resp := utils.NewAPIResponse("", nil)
	var device *models.Device
	err := decoder.Decode(&device)
	if err != nil {
		resp.Message = "Invalid JSON"
		d.e.Log.WithField("Err", err).Error("Invalid JSON")
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if r.Method == "PUT" && device.Slug != p.ByName("slug") {
		resp.Message = "Slug mismatch"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	device.SetEnv(d.e)

	err = device.Save()
	if err != nil {
		resp.Message = err.Error()
		d.e.Log.WithField("Err", err).Error("Failed to save device")
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	resp.Message = "Device saved successfully"
	resp.Data = device
	resp.WriteResponse(w, http.StatusOK)
}
