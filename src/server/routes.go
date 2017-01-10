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
	r.GET("/devices", d.ShowDeviceList)

	c := controllers.GetConfigController(e)
	r.GET("/devices/:slug", c.ShowDeviceConfigList)
	r.GET("/devices/:slug/:config", c.ShowConfig)
	r.GET("/configs/:config", c.ShowConfig)

	t := controllers.GetConnProfileController(e)
	r.GET("/profiles", t.ShowTypeList)
	r.GET("/profiles/:slug", t.ShowType)

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

	t := controllers.GetConnProfileController(e)
	r.GET("/api/profiles", t.ApiGetTypes)
	r.GET("/api/profiles/:slug", t.ApiGetTypes)

	l := controllers.GetLogController(e)
	r.GET("/api/logs", l.ApiGetLogs)

	return r
}

func apiPUTRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.GetDeviceController(e)
	r.PUT("/api/devices/:slug", d.ApiPutDevice)

	t := controllers.GetConnProfileController(e)
	r.PUT("/api/profiles/:slug", t.ApiPutType)

	return r
}

func apiPOSTRoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.GetDeviceController(e)
	r.POST("/api/devices", d.ApiPostDevice)

	j := controllers.GetJobController(e)
	r.POST("/api/job/start", j.ApiStartJob)
	r.POST("/api/job/stop/:id", j.ApiStopJob)

	t := controllers.GetConnProfileController(e)
	r.POST("/api/profiles", t.ApiPostType)

	return r
}

func apiDELETERoutes(e *utils.Environment) http.Handler {
	r := httprouter.New()

	d := controllers.GetDeviceController(e)
	r.DELETE("/api/devices/:slug", d.ApiDeleteDevice)

	t := controllers.GetConnProfileController(e)
	r.DELETE("/api/profiles/:slug", t.ApiDeleteType)

	return r
}
