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
	r.GET("/devices/:id", d.ShowDevice)
	r.GET("/devices/:id/:config", d.ShowConfig)

	m := controllers.NewManager(e)
	r.GET("/manage/:object", m.Manage)

	r.Handler("GET", "/api/*a", apiGETRoutes(e))
	r.Handler("PUT", "/api/*a", apiRoutes(e))
	r.Handler("POST", "/api/*a", apiRoutes(e))
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
		w.Write([]byte(`"API not yet"`))
	})
}

func apiGETRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.NewDevice(e)
	r.GET("/api/devices", d.ApiGetDevices)
	r.GET("/api/devices/:id", d.ApiGetDevices)
	r.GET("/api/devices/:id/configs", d.ApiGetConfigs)
	r.GET("/api/devices/:id/configs/:config", d.ApiGetConfigs)

	return r
}
