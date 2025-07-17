package logging

import (
	"github.com/agentic-layer/model-router-krakend/lib/header"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type PluginLogger struct {
	logger    *log.Logger
	sessionId string
}

var zero = time.Time{}

func newLogger(pluginName string) *log.Logger {
	if testing.Testing() && len(pluginName) > 14 {
		panic("pluginName '" + pluginName + "' > 14 characters")
	}
	return log.New(os.Stdout, fmt.Sprintf("[%-14s] ", strings.ToTitle(pluginName)), log.LstdFlags|log.Lmicroseconds)
}

func newRandomId() string {
	return uuid.New().String()
}

// New creates a new default logger with no sessionId
func New(pluginName string) *PluginLogger {
	return &PluginLogger{
		logger:    newLogger(pluginName),
		sessionId: "",
	}
}

// NewWithRandomId creates a new logger with a random sessionId
func NewWithRandomId(pluginName string) *PluginLogger {
	return &PluginLogger{
		logger:    newLogger(pluginName),
		sessionId: newRandomId(),
	}
}

// NewFromHttpRequest creates a new logger using the sessionId from the request
func NewFromHttpRequest(pluginName string, req *http.Request) *PluginLogger {
	sessionId := req.Header.Get(header.SessionId)

	return &PluginLogger{
		logger:    newLogger(pluginName),
		sessionId: sessionId,
	}
}

func (l *PluginLogger) output(level string, start time.Time, format string, v ...interface{}) {
	timer := ""
	if !start.IsZero() {
		timer = fmt.Sprintf(" (took %v)", time.Since(start))
	}
	if l.sessionId == "" {
		l.logger.Printf("| %-5s | %v%s\n", level, fmt.Sprintf(format, v...), timer)
	} else {
		l.logger.Printf("| %s | %-5s | %v%s\n", l.sessionId, level, fmt.Sprintf(format, v...), timer)
	}
}

func (l *PluginLogger) Debug(format string, v ...interface{}) {
	l.output("DEBUG", zero, format, v...)
}

func (l *PluginLogger) DebugTimed(start time.Time, format string, v ...interface{}) {
	l.output("DEBUG", start, format, v...)
}

func (l *PluginLogger) Info(format string, v ...interface{}) {
	l.output("INFO", zero, format, v...)
}

func (l *PluginLogger) InfoTimed(start time.Time, format string, v ...interface{}) {
	l.output("INFO", start, format, v...)
}

func (l *PluginLogger) Warn(format string, v ...interface{}) {
	l.output("WARN", zero, format, v...)
}

func (l *PluginLogger) WarnTimed(start time.Time, format string, v ...interface{}) {
	l.output("WARN", start, format, v...)
}

func (l *PluginLogger) Error(format string, v ...interface{}) {
	l.output("ERROR", zero, format, v...)
}

func (l *PluginLogger) ErrorTimed(start time.Time, format string, v ...interface{}) {
	l.output("ERROR", start, format, v...)
}
