package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
)

var deviceControllerSingle *DeviceController

type DeviceController struct {
	e *utils.Environment
}

func newDeviceController(e *utils.Environment) *DeviceController {
	return &DeviceController{e: e}
}

func GetDeviceController(e *utils.Environment) *DeviceController {
	if deviceControllerSingle == nil {
		deviceControllerSingle = newDeviceController(e)
	}
	return deviceControllerSingle
}

func (d *DeviceController) ShowDevice(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if p.ByName("slug") != "" {
		d.showDeviceConfigList(w, r, p.ByName("slug"))
		return
	}

	devices, err := models.GetAllDevices(d.e)
	if err != nil {
		d.e.Log.WithField("Err", err).Error("Failed to get devices")
		return
	}

	data := map[string]interface{}{
		"devices": devices,
	}
	d.e.View.NewView("device-list", r).Render(w, data)
	return
}

func (d *DeviceController) showDeviceConfigList(w http.ResponseWriter, r *http.Request, slug string) {
	device, err := models.GetDeviceBySlug(d.e, slug)
	if err != nil {
		d.e.Log.WithField("Err", err).Error("Failed to get device")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if device.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	configs, err := models.GetConfigsForDevice(d.e, device.ID)
	if err != nil {
		d.e.Log.WithField("Err", err).Error("Failed to get configs")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"device":  device,
		"configs": configs,
	}
	d.e.View.NewView("device", r).Render(w, data)
}

func (d *DeviceController) ShowConfig(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("slug")
	configSlug := p.ByName("config")

	if configSlug == "" {
		d.showDeviceConfigList(w, r, name)
		return
	}

	device, err := models.GetDeviceBySlug(d.e, name)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get device")
		return
	}

	if device.ID == 0 {
		d.e.Log.Debug("Device not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	config, err := models.GetConfigBySlug(d.e, configSlug)
	if err != nil {
		d.e.Log.WithField("error", err).Error("Couldn't get config")
		return
	}

	if config.ID == 0 {
		d.e.Log.Debug("Config not found")
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

func (d *DeviceController) ApiGetDevices(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func (d *DeviceController) ApiPostDevice(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	resp := utils.NewAPIResponse("", nil)
	var apiDevice *models.Device
	err := decoder.Decode(&apiDevice)
	if err != nil {
		resp.Message = "Invalid JSON"
		d.e.Log.WithField("Err", err).Error("Invalid JSON")
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	postedDevice := models.NewDevice(d.e)

	postedDevice.Profile, err = models.GetTypeByID(d.e, apiDevice.Profile.ID)
	if err != nil {
		resp.Message = "Unknown device type"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	if apiDevice.Address == "" ||
		apiDevice.Name == "" {
		resp.Message = "Missing data field"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	postedDevice.Address = apiDevice.Address
	postedDevice.Name = apiDevice.Name

	if err := postedDevice.Save(); err != nil {
		resp.Message = "Failed to save device"
		d.e.Log.WithField("Err", err).Error("Failed to save device")
		resp.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	resp.Message = "Device saved successfully"
	resp.Data = postedDevice
	resp.WriteResponse(w, http.StatusOK)
}

func (d *DeviceController) ApiPutDevice(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	resp := utils.NewAPIResponse("", nil)
	var apiDevice *models.Device
	err := decoder.Decode(&apiDevice)
	if err != nil {
		resp.Message = "Invalid JSON"
		d.e.Log.WithField("Err", err).Error("Invalid JSON")
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Retrieve device from database and work on it
	originalDevice, err := models.GetDeviceBySlug(d.e, p.ByName("slug"))
	if err != nil {
		resp.Message = "Device " + p.ByName("slug") + " was not found"
		resp.WriteResponse(w, http.StatusBadRequest)
		return
	}

	// Copy data from api request if different
	if apiDevice.Profile.ID > 0 && apiDevice.Profile.ID != originalDevice.Profile.ID {
		originalDevice.Profile, err = models.GetTypeByID(d.e, apiDevice.Profile.ID)
		if err != nil {
			resp.Message = "Unknown device type"
			resp.WriteResponse(w, http.StatusBadRequest)
			return
		}
	}

	if apiDevice.Name != "" {
		originalDevice.Name = apiDevice.Name
	}

	if apiDevice.Address != "" {
		originalDevice.Address = apiDevice.Address
	}

	if err := originalDevice.Save(); err != nil {
		resp.Message = "Error saving device"
		d.e.Log.WithField("Err", err).Error("Failed to save device")
		resp.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	resp.Message = "Device saved successfully"
	resp.Data = originalDevice
	resp.WriteResponse(w, http.StatusOK)
}

func (d *DeviceController) ApiDeleteDevice(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("slug")

	ret := utils.NewAPIResponse("", nil)
	if name == "" {
		ret.Message = "No device given"
		ret.WriteResponse(w, http.StatusBadRequest)
		return
	}

	device, err := models.GetDeviceBySlug(d.e, name)
	if err != nil {
		ret.Message = "Error getting device"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	if device.ID == 0 { // No device with that slug, return
		ret.WriteResponse(w, http.StatusNoContent)
		return
	}

	if err := device.Delete(); err != nil {
		ret.Message = err.Error()
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	configs, err := models.GetConfigsForDevice(d.e, device.ID)
	if err != nil {
		ret.Message = "Error deleting configs"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	finishedWithErrors := false
	for _, config := range configs {
		if err := config.Delete(); err != nil {
			finishedWithErrors = true
		}
	}

	if finishedWithErrors {
		ret.Message = "Device deleted, but configurations were not"
		ret.WriteResponse(w, http.StatusInternalServerError)
		return
	}

	ret.WriteResponse(w, http.StatusNoContent)
	return
}
