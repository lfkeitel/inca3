// This source file is part of the Inca project.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

// Config is a project-wide struct that holds configuration information
type Config struct {
	sourceFile string
	Core       struct {
		SiteTitle          string
		SiteCompanyName    string
		SiteDomainName     string
		SiteFooterText     string
		JobSchedulerWakeUp string
	}
	Logging struct {
		Enabled    bool
		EnableHTTP bool
		Level      string
		Path       string
	}
	Database struct {
		Path string
	}
	Webserver struct {
		Address             string
		HTTPPort            int
		HTTPSPort           int
		TLSCertFile         string
		TLSKeyFile          string
		RedirectHTTPToHTTPS bool
	}
	Configs struct {
		BaseDir string
	}
}

// FindConfigFile will locate the a configuration file by looking at the following places
// and choosing the first: INCA_CONFIG env variable, $PWD/config.toml, $PWD/config/config.toml,
// $HOME/.inca/config.toml, and /etc/inca/config.toml.
func FindConfigFile() string {
	if os.Getenv("INCA_CONFIG") != "" && FileExists(os.Getenv("INCA_CONFIG")) {
		return os.Getenv("INCA_CONFIG")
	} else if FileExists("./config.toml") {
		return "./config.toml"
	} else if FileExists("./config/config.toml") {
		return "./config/config.toml"
	} else if FileExists(os.ExpandEnv("$HOME/.inca/config.toml")) {
		return os.ExpandEnv("$HOME/.inca/config.toml")
	} else if FileExists("/etc/inca/config.toml") {
		return "/etc/inca/config.toml"
	}
	return ""
}

// NewConfig will load a Config object using the given TOML file.
func NewConfig(configFile string) (conf *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	if configFile == "" {
		configFile = "config.toml"
	}

	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var con Config
	if err := toml.Unmarshal(buf, &con); err != nil {
		return nil, err
	}
	con.sourceFile = configFile
	return &con, nil
}
