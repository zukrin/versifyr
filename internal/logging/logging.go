package logging

import (
	"github.com/pterm/pterm"
)

type LogLevel int

const (
	All LogLevel = iota
	Debug
	Info
	Error
)

func (s LogLevel) String() string {
	switch s {
	case All:
		return "all"
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Error:
		return "error"
	}
	return "unknown"
}

func LogLevelFromString(s string) LogLevel {
	switch s {
	case "all":
		return All
	case "debug":
		return Debug
	case "info":
		return Info
	case "error":
		return Error
	}
	return All
}

type Logger struct {
	InfoWriter  *infoWriter
	ErrWriter   *errWriter
	DebugWriter *debugWriter
	LogLevel    LogLevel
}

func NewLogger() *Logger {
	l := &Logger{}
	l.InfoWriter = &infoWriter{l}
	l.ErrWriter = &errWriter{l}
	l.DebugWriter = &debugWriter{l}
	l.LogLevel = Info
	pterm.ThemeDefault.DebugMessageStyle = *pterm.NewStyle(pterm.FgLightWhite)
	return l
}

func (l *Logger) Debug(template string, args ...interface{}) {
	if Debug >= l.LogLevel {
		pterm.ThemeDefault.DebugMessageStyle.Printfln(template, args...)
	}
}

func (l *Logger) Info(template string, args ...interface{}) {
	if Info >= l.LogLevel {
		pterm.ThemeDefault.InfoMessageStyle.Printfln(template, args...)
	}
}

func (l *Logger) Error(template string, args ...interface{}) {
	if Error >= l.LogLevel {
		pterm.ThemeDefault.ErrorMessageStyle.Printfln(template, args...)
	}
}

type debugWriter struct {
	Logger *Logger
}

func (dw *debugWriter) Write(p []byte) (n int, err error) {
	pterm.ThemeDefault.DebugMessageStyle.Printf(string(p))
	return len(p), nil
}

type infoWriter struct {
	Logger *Logger
}

func (iw *infoWriter) Write(p []byte) (n int, err error) {
	pterm.ThemeDefault.InfoMessageStyle.Printf(string(p))
	return len(p), nil
}

type errWriter struct {
	Logger *Logger
}

func (iw *errWriter) Write(p []byte) (n int, err error) {
	pterm.ThemeDefault.ErrorMessageStyle.Printf(string(p))
	return len(p), nil
}
