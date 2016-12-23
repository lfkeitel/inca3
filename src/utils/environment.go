// This source file is part of the Inca project.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

// EnvType is the runtime state of the application
type EnvType string

const (
	// EnvTesting is used when running automated tests
	EnvTesting EnvType = "testing"
	// EnvProd is used during normal operation
	EnvProd EnvType = "production"
	// EnvDev is used during development
	EnvDev EnvType = "development"
)

// An Environment is a project-wide struct that holds resources needed by all
// parts of the application.
type Environment struct {
	Config *Config
	DB     *DatabaseAccessor
	Log    *Logger
	Env    EnvType
	View   *Views
}

// NewEnvironment will create a new Environment object using run type e.
func NewEnvironment(e EnvType) *Environment {
	return &Environment{Env: e}
}

func (e *Environment) IsDev() bool {
	return e.Env == EnvDev
}
