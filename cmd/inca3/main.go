// This source file is part of the Inca project.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lfkeitel/inca3/src/server"
	"github.com/lfkeitel/inca3/src/utils"
	"github.com/lfkeitel/verbose"
)

var (
	configFile string
	dev        bool
	verFlag    bool
	testConfig bool

	version   = ""
	buildTime = ""
	builder   = ""
	goversion = ""
)

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration file path")
	flag.BoolVar(&dev, "d", false, "Run in development mode")
	flag.BoolVar(&verFlag, "version", false, "Display version information")
	flag.BoolVar(&verFlag, "v", verFlag, "Display version information")
	flag.BoolVar(&testConfig, "t", false, "Test main configuration")
}

func main() {
	flag.Parse()

	if verFlag {
		displayVersionInfo()
		return
	}

	if configFile == "" || !utils.FileExists(configFile) {
		configFile = utils.FindConfigFile()
	}
	if configFile == "" {
		fmt.Println("No configuration file found")
		os.Exit(1)
	}

	if testConfig {
		testMainConfig()
		return
	}

	var err error
	e := utils.NewEnvironment(utils.EnvProd)
	if dev {
		e.Env = utils.EnvDev
	}

	e.Config, err = utils.NewConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		os.Exit(1)
	}

	e.Log = utils.NewLogger(e.Config, "inca")
	e.Log.WithFields(verbose.Fields{
		"path": configFile,
	}).Info("Loaded configuration")

	e.DB, err = utils.NewDatabaseAccessor(e.Config)
	if err != nil {
		e.Log.WithField("error", err).Fatal("Error loading database")
	}
	e.Log.WithFields(verbose.Fields{
		"path": e.Config.Database.Path,
	}).Debug("Loaded database")

	e.View, err = utils.NewViewer(e, "public/templates")
	if err != nil {
		e.Log.WithField("error", err).Fatal("Error loading frontend templates")
	}

	server.NewServer(e, server.LoadRoutes(e)).Run()
}

func displayVersionInfo() {
	fmt.Printf(`Inca3 - (C) 2016 University of Southern Indiana

Version:     %s
Built:       %s
Compiled by: %s
Go version:  %s
`, version, buildTime, builder, goversion)
}

func testMainConfig() {
	_, err := utils.NewConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Configuration looks good")
}
