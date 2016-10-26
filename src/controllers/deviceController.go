package controllers

import (
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
	name := p.ByName("id")

	if name != "" {
		d.showDeviceConfigList(w, r, name)
		return
	}

	// Show all devices
	devices, err := models.GetAllDevices(d.e)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get devices")
		return
	}

	data := map[string]interface{}{
		"devices": devices,
	}
	d.e.View.NewView("device-list", r).Render(w, data)
	return
}

func (d *Device) showDeviceConfigList(w http.ResponseWriter, r *http.Request, name string) {
	device, err := models.GetDeviceByID(d.e, name)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get device")
		return
	}

	if device.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	configs, err := models.GetConfigsForDevice(d.e, device.ID)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get configs")
		return
	}

	data := map[string]interface{}{
		"device":  device,
		"configs": configs,
	}
	d.e.View.NewView("device", r).Render(w, data)
}

func (d *Device) ShowConfig(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("id")
	configID := p.ByName("config")

	if configID == "" {
		d.showDeviceConfigList(w, r, name)
		return
	}

	device, err := models.GetDeviceByID(d.e, name)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get device")
		return
	}

	if device.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	config, err := models.GetConfigByID(d.e, configID)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get config")
		return
	}

	if config.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	configText, err := config.GetText()
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get config text")
		return
	}

	data := map[string]interface{}{
		"device":     device,
		"config":     config,
		"configText": string(configText),
	}
	d.e.View.NewView("config", r).Render(w, data)
}

func (d *Device) ApiGetDevices(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("id")

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

	device, err := models.GetDeviceByID(d.e, name)
	if err != nil {
		ret.Message = "Error getting devices"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	ret.Data = device
	ret.WriteResponse(w, http.StatusOK)
	return
}

func (d *Device) ApiGetConfigs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("id")
	configID := p.ByName("config")

	ret := utils.NewAPIResponse("", nil)
	if name == "" {
		ret.WriteResponse(w, http.StatusNotFound)
		return
	}

	// Check device exists
	device, err := models.GetDeviceByID(d.e, name)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get device")
		return
	}

	if device.ID == "" {
		ret.WriteResponse(w, http.StatusNotFound)
		return
	}

	if configID == "" { // Return all configs
		configs, err := models.GetConfigsForDevice(d.e, device.ID)
		if err != nil {
			d.e.Log.WithField("error", err).Error("Couldn't get configs")
			return
		}

		ret.Data = configs
		ret.WriteResponse(w, http.StatusOK)
		return
	}

	config, err := models.GetConfigByID(d.e, configID)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get configs")
		return
	}

	ret.Data = config
	ret.WriteResponse(w, http.StatusOK)
	return
}
