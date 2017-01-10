package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/lfkeitel/verbose"
)

var strToLogLevels = map[string]verbose.LogLevel{
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

// type Level verbose.LogLevel

// var (
// 	Debug     Level = verbose.LogLevelDebug
// 	Info      Level = verbose.LogLevelInfo
// 	Notice    Level = verbose.LogLevelNotice
// 	Warning   Level = verbose.LogLevelWarning
// 	Error     Level = verbose.LogLevelError
// 	Critical  Level = verbose.LogLevelCritical
// 	Alert     Level = verbose.LogLevelAlert
// 	Emergency Level = verbose.LogLevelEmergency
// 	Fatal     Level = verbose.LogLevelFatal
// )

var SystemLogger *Logger

func init() {
	SystemLogger = NewEmptyLogger()
}

type Log struct {
	Level     verbose.LogLevel `json:"level"`
	Message   string           `json:"message"`
	Timestamp int64            `json:"timestamp"`
}

type Logger struct {
	*verbose.Logger
	c           *Config
	numUserLogs int
	userLogs    []*Log
}

func NewEmptyLogger() *Logger {
	return &Logger{
		Logger:      verbose.New("null"),
		c:           &Config{},
		numUserLogs: 0,
		userLogs:    make([]*Log, 20),
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

	if level, ok := strToLogLevels[strings.ToLower(c.Logging.Level)]; ok {
		sh.SetMinLevel(level)
		fh.SetMinLevel(level)
	}
	fh.SetFormatter(&verbose.JSONFormatter{})
	return &Logger{
		Logger:      logger,
		c:           c,
		numUserLogs: 0,
		userLogs:    make([]*Log, 20),
	}
}

func (l *Logger) UserLog(lvl verbose.LogLevel, f string, args ...interface{}) {
	// Shift all logs down by one
	for i := len(l.userLogs) - 1; i > 0; i-- {
		l.userLogs[i] = l.userLogs[i-1]
	}

	if l.numUserLogs < len(l.userLogs) {
		l.numUserLogs++
	}

	l.userLogs[0] = &Log{
		Level:     lvl,
		Message:   fmt.Sprintf(f, args...),
		Timestamp: time.Now().Unix(),
	}
}

func (l *Logger) GetUserLogs() []*Log {
	return l.userLogs[:l.numUserLogs]
}
