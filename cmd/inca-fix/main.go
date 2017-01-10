// This source file is part of the Inca project.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lfkeitel/inca3/src/models"
	"github.com/lfkeitel/inca3/src/utils"
	"github.com/lfkeitel/verbose"
)

var (
	e *utils.Environment

	configFile string
	dev        bool
	verFlag    bool

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

	var err error
	e = utils.NewEnvironment(utils.EnvProd)
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

	fixConfigs()
}

func displayVersionInfo() {
	fmt.Printf(`Inca3 - (C) 2016 University of Southern Indiana

Version:     %s
Built:       %s
Compiled by: %s
Go version:  %s
`, version, buildTime, builder, goversion)
}

func fixConfigs() {
	dirs, err := ioutil.ReadDir(e.Config.DirPaths.BaseDir)
	if err != nil {
		panic(err)
	}

	count := 0

	for _, dir := range dirs {
		device, err := models.GetDeviceByIP(e, dir.Name())
		if err != nil || device == nil {
			e.Log.Errorf("Failed to load device %s", dir.Name())
			continue
		}

		configs, err := ioutil.ReadDir(filepath.Join(e.Config.DirPaths.BaseDir, dir.Name()))
		if err != nil {
			panic(err)
		}

		for _, config := range configs {
			c, err := models.GetConfigByFilename(e, filepath.Join(dir.Name(), config.Name()))
			if err != nil {
				e.Log.Errorf("Failed to load config %s", config.Name())
				continue
			}

			if c != nil { // config is already in the database
				continue
			}

			date := strings.SplitN(config.Name(), ".", 2)[0]

			c = models.NewConfig(e)
			c.Slug = device.Slug + "_" + date
			c.DeviceID = device.ID
			c.Filename = filepath.Join(dir.Name(), config.Name())
			c.Created = config.ModTime()
			c.Compressed = true

			if err := c.Save(); err != nil {
				e.Log.WithField("Err", err).Error("Failed to save config")
				continue
			}

			e.Log.Infof("Added log for %s from %s", device.Slug, date)
			count++
		}
	}

	e.Log.Infof("Added %d configs", count)
}
