package utils

import (
	"strings"
	"time"

	"github.com/lfkeitel/verbose"
)

var logLevels = map[string]verbose.LogLevel{
	"debug":     verbose.LogLevelDebug,
	"info":      verbose.LogLevelInfo,
	"notice":    verbose.LogLevelNotice,
	"warning":   verbose.LogLevelWarning,
	"error":     verbose.LogLevelError,
	"critical":  verbose.LogLevelCritical,
	"alert":     verbose.LogLevelAlert,
	"emergency": verbose.LogLevelEmergency,
	"fatal":     verbose.LogLevelFatal,
}

var SystemLogger *Logger

func init() {
	SystemLogger = NewEmptyLogger()
}

type Logger struct {
	*verbose.Logger
	c      *Config
	timers map[string]time.Time
}

func NewEmptyLogger() *Logger {
	return &Logger{
		Logger: verbose.New("null"),
		timers: make(map[string]time.Time),
		c:      &Config{},
	}
}

func NewLogger(c *Config, name string) *Logger {
	logger := verbose.New(name)
	if !c.Logging.Enabled {
		return &Logger{
			Logger: logger,
		}
	}
	sh := verbose.NewStdoutHandler(true)
	fh, _ := verbose.NewFileHandler(c.Logging.Path)
	logger.AddHandler("stdout", sh)
	logger.AddHandler("file", fh)

	if level, ok := logLevels[strings.ToLower(c.Logging.Level)]; ok {
		sh.SetMinLevel(level)
		fh.SetMinLevel(level)
	}
	fh.SetFormatter(&verbose.JSONFormatter{})
	return &Logger{
		Logger: logger,
		c:      c,
	}
}
