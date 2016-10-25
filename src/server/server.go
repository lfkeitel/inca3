// This source file is part of the Inca project.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lfkeitel/verbose"

	"github.com/lfkeitel/inca3/src/utils"

	"gopkg.in/tylerb/graceful.v1"
)

type Server struct {
	e         *utils.Environment
	routes    http.Handler
	address   string
	httpPort  string
	httpsPort string
}

func NewServer(e *utils.Environment, routes http.Handler) *Server {
	serv := &Server{
		e:       e,
		routes:  routes,
		address: e.Config.Webserver.Address,
	}

	serv.httpPort = strconv.Itoa(e.Config.Webserver.HTTPPort)
	serv.httpsPort = strconv.Itoa(e.Config.Webserver.HTTPSPort)
	return serv
}

func (s *Server) Run() {
	if s.e.Config.Webserver.TLSCertFile == "" || s.e.Config.Webserver.TLSKeyFile == "" {
		s.startHttp()
		return
	}

	if s.e.Config.Webserver.RedirectHTTPToHTTPS {
		go s.startRedirector()
	}
	s.startHttps()
}

func (s *Server) startRedirector() {
	s.e.Log.WithFields(verbose.Fields{
		"address": s.address,
		"port":    s.httpPort,
	}).Debug("Starting HTTP->HTTPS redirector")
	timeout := 1 * time.Second
	if s.e.IsDev() {
		timeout = 1 * time.Millisecond
	}
	srv := &graceful.Server{
		Timeout: timeout,
		Server:  &http.Server{Addr: s.address + ":" + s.httpPort, Handler: http.HandlerFunc(s.redirectToHttps)},
	}
	if err := srv.ListenAndServe(); err != nil {
		s.e.Log.Fatal(err)
	}
}

func (s *Server) startHttp() {
	s.e.Log.WithFields(verbose.Fields{
		"address": s.address,
		"port":    s.httpPort,
	}).Info("Starting HTTP server")
	timeout := 5 * time.Second
	if s.e.IsDev() {
		timeout = 1 * time.Millisecond
	}
	srv := &graceful.Server{
		Timeout: timeout,
		Server:  &http.Server{Addr: s.address + ":" + s.httpPort, Handler: s.routes},
	}
	if err := srv.ListenAndServe(); err != nil {
		s.e.Log.Fatal(err)
	}
}

func (s *Server) startHttps() {
	s.e.Log.WithFields(verbose.Fields{
		"address": s.address,
		"port":    s.httpsPort,
	}).Info("Starting HTTPS server")
	timeout := 5 * time.Second
	if s.e.IsDev() {
		timeout = 1 * time.Millisecond
	}
	srv := &graceful.Server{
		Timeout: timeout,
		Server:  &http.Server{Addr: s.address + ":" + s.httpsPort, Handler: s.routes},
	}
	if err := srv.ListenAndServeTLS(
		s.e.Config.Webserver.TLSCertFile,
		s.e.Config.Webserver.TLSKeyFile,
	); err != nil {
		s.e.Log.Fatal(err)
	}
}

func (s *Server) redirectToHttps(w http.ResponseWriter, r *http.Request) {
	// Lets not do a split if we don't need to
	if s.httpPort == "80" && s.httpsPort == "443" {
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		return
	}

	host := strings.Split(r.Host, ":")[0]
	http.Redirect(w, r, "https://"+host+":"+s.httpsPort+r.RequestURI, http.StatusMovedPermanently)
}
