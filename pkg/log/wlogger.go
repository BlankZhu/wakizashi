// Package wlogger describe the liveness probe used by wakizashi.
// ATTENTION: this package is deprecated
// This is possible for k8s or other custom monitor to check wakizashi's liveness.
// Example:
// 	wlog := wlogger.Get()
// 	wlog.SetLevel(wlogger.InfoLevel)
// 	...
// 	wlog.Fatalf("error exiting... detail:%s", err.Error())
package wlogger

import (
	"log"
	"os"
	"sync"
)

const (
	// DebugLevel debug log level
	DebugLevel = iota
	// InfoLevel info log level
	InfoLevel
	// WarningLevel warning log level
	WarningLevel
	// ErrorLevel error log level
	ErrorLevel
	// FatalLevel fatal log level
	FatalLevel
)

type wlogger struct {
	logLevel int
	logger   *log.Logger
}

var wlog *wlogger
var once sync.Once

// Get return the singleton wlogger
func Get() *wlogger {
	once.Do(func() {
		wlog = &wlogger{}
		wlog.init()
	})
	return wlog
}

// Init initalize the wakizashi's logger
func (l *wlogger) init() {
	l.logLevel = DebugLevel
	l.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

// SwitchLogLevel switch the logging level to the value between: [DebugLevel, FatalLevel]
func (l *wlogger) SetLevel(logLevel int) {
	if logLevel >= DebugLevel && logLevel <= FatalLevel {
		l.logLevel = logLevel
	}
}

// Debug log given string on debug level
func (l *wlogger) Debug(str string) {
	if l.logLevel <= DebugLevel {
		l.logger.Printf("[DEUBG] %s", str)
	}
}

// Debugf log given formatted string with parameters on debug level
func (l *wlogger) Debugf(format string, any ...interface{}) {
	if l.logLevel <= DebugLevel {
		l.logger.Printf("[DEBUG] "+format, any...)
	}
}

// Info log given string on info level
func (l *wlogger) Info(str string) {
	if l.logLevel <= InfoLevel {
		l.logger.Printf("[INFO] %s", str)
	}
}

// Infof log given formatted string with parameters on info level
func (l *wlogger) Infof(format string, any ...interface{}) {
	if l.logLevel <= InfoLevel {
		l.logger.Printf("[INFO] "+format, any...)
	}
}

// Warning log given string on warning level
func (l *wlogger) Warning(str string) {
	if l.logLevel <= WarningLevel {
		l.logger.Printf("[WARNING] %s", str)
	}
}

// Warningf log given formatted string with parameters on warning level
func (l *wlogger) Warningf(format string, any ...interface{}) {
	if l.logLevel <= WarningLevel {
		l.logger.Printf("[WARNING] "+format, any...)
	}
}

// Error log given string on error level
func (l *wlogger) Error(str string) {
	if l.logLevel <= ErrorLevel {
		l.logger.Printf("[ERROR] %s", str)
	}
}

// Errorf log given formatted string with parameters on error level
func (l *wlogger) Errorf(format string, any ...interface{}) {
	if l.logLevel <= ErrorLevel {
		l.logger.Printf("[ERROR] "+format, any...)
	}
}

// Fatal log given string on fatal level and exit the wakizashi
func (l *wlogger) Fatal(str string) {
	if l.logLevel <= FatalLevel {
		l.logger.Printf("[FATAL] %s", str)
	}
	os.Exit(1)
}

// Fatalf log given formatted string with parameters on fatal level and exit the wakizashi
func (l *wlogger) Fatalf(format string, any ...interface{}) {
	if l.logLevel <= FatalLevel {
		l.logger.Printf("[FATAL] "+format, any...)
	}
	os.Exit(1)
}
