// This source file is part of the Inca project.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lfkeitel/inca3/src/utils"
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

	if configFile == "" || !utils.FileExists(configFile) {
		configFile = utils.FindConfigFile()
	}
	if configFile == "" {
		fmt.Println("No configuration file found")
		os.Exit(1)
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

	e.Log = utils.NewLogger(e.Config
	e.Log.Debugf("Configuration loaded from %s", configFile)


	e.DB, err = utils.NewDatabaseAccessor(e.Config)
	if err != nil {
		e.Log.WithField("error", err).Fatal("Error loading database")
	}
	e.Log.WithFields(verbose.Fields{
		"type":    e.Config.Database.Type,
		"address": e.Config.Database.Address,
	}).Debug("Loaded database")

	e.View, err = utils.NewViewer(e, "public/templates")
	if err != nil {
		e.Log.WithField("error", err).Fatal("Error loading frontend templates")
	}
	// Start server
}
