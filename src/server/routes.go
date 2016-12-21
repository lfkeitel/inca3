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

	d := controllers.GetDeviceController(e)
	r.GET("/devices", d.ShowDevice)
	r.GET("/devices/:slug", d.ShowDevice)
	r.GET("/devices/:slug/:config", d.ShowConfig)

	r.Handler("GET", "/api/*a", apiGETRoutes(e))
	r.Handler("PUT", "/api/*a", apiPUTRoutes(e))
	r.Handler("POST", "/api/*a", apiPOSTRoutes(e))
	r.Handler("DELETE", "/api/*a", apiDELETERoutes(e))
	return r
}

func rootHandler(e *utils.Environment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.View.NewView("index", r).Render(w, nil)
	})
}

func apiGETRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.GetDeviceController(e)
	r.GET("/api/devices", d.ApiGetDevices)
	r.GET("/api/devices/:slug", d.ApiGetDevices)

	c := controllers.GetConfigController(e)
	r.GET("/api/devices/:slug/configs", c.ApiGetDeviceConfigs)
	r.GET("/api/devices/:slug/configs/:config", c.ApiGetDeviceConfigs)
	r.GET("/api/configs/:config", c.ApiGetConfig)

	j := controllers.GetJobController(e)
	r.GET("/api/job/status/:id", j.ApiJobStatus)

	t := controllers.GetTypeController(e)
	r.GET("/api/types", t.ApiGetTypes)
	r.GET("/api/types/:id", t.ApiGetTypes)

	return r
}

func apiPUTRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.GetDeviceController(e)
	r.PUT("/api/devices/:slug", d.ApiPutDevice)

	t := controllers.GetTypeController(e)
	r.PUT("/api/types/:id", t.ApiPutType)

	return r
}

func apiPOSTRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.GetDeviceController(e)
	r.POST("/api/devices", d.ApiPostDevice)

	j := controllers.GetJobController(e)
	r.POST("/api/job/start", j.ApiStartJob)
	r.POST("/api/job/stop/:id", j.ApiStopJob)

	t := controllers.GetTypeController(e)
	r.POST("/api/types/:id", t.ApiPostType)

	return r
}

func apiDELETERoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.GetDeviceController(e)
	r.DELETE("/api/devices/:slug", d.ApiDeleteDevice)

	t := controllers.GetTypeController(e)
	r.DELETE("/api/types/:id", t.ApiDeleteType)

	return r
}
