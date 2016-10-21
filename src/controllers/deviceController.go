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

	// Show all devices
	if name == "" {
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

	// Show a particular
	device, err := models.GetDeviceByID(d.e, name)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get devices")
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
	// If config is empty, show a page with all the device's known configs
	// If config is not empty, show the config
}
