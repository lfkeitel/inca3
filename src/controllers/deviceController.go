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
	if p.ByName("id") != "" {
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
	name := p.ByName("id")
	configID := p.ByName("config")

	if configID == "" {
		d.showDeviceConfigList(w, r)
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
