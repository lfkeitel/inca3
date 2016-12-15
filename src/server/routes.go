package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lfkeitel/inca3/src/controllers"
	"github.com/lfkeitel/inca3/src/utils"
)

func LoadRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()
	r.ServeFiles("/static/*filepath", http.Dir("./public/static"))
	r.Handler("GET", "/", rootHandler(e))

	d := controllers.NewDevice(e)
	r.GET("/devices", d.ShowDevice)
	r.GET("/devices/:slug", d.ShowDevice)
	r.GET("/devices/:slug/:config", d.ShowConfig)

	m := controllers.NewManager(e)
	r.GET("/manage/:object", m.Manage)

	r.Handler("GET", "/api/*a", apiGETRoutes(e))
	r.Handler("PUT", "/api/*a", apiPUTRoutes(e))
	r.Handler("POST", "/api/*a", apiPOSTRoutes(e))
	r.Handler("DELETE", "/api/*a", apiRoutes(e))
	return r
}

func rootHandler(e *utils.Environment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.View.NewView("index", r).Render(w, nil)
	})
}

func apiRoutes(e *utils.Environment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.NewAPIResponse("API not implemented", nil).WriteResponse(w, http.StatusOK)
	})
}

func apiGETRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.NewDevice(e)
	r.GET("/api/devices", d.ApiGetDevices)
	r.GET("/api/devices/:slug", d.ApiGetDevices)

	c := controllers.NewConfig(e)
	r.GET("/api/devices/:slug/configs", c.ApiGetDeviceConfigs)
	r.GET("/api/devices/:slug/configs/:config", c.ApiGetDeviceConfigs)
	r.GET("/api/configs/:config", c.ApiGetConfig)

	return r
}

func apiPUTRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.NewDevice(e)
	r.PUT("/api/devices/:slug", d.ApiSaveDevice)

	return r
}

func apiPOSTRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.NewDevice(e)
	r.POST("/api/devices", d.ApiSaveDevice)

	return r
}
