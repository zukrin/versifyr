package logging

import (
	"github.com/pterm/pterm"
)

type Logger struct {
	InfoWriter *infoWriter
	ErrWriter  *errWriter
}

func NewLogger() *Logger {
	l := &Logger{}
	l.InfoWriter = &infoWriter{l}
	l.ErrWriter = &errWriter{l}
	return l
}

func (l *Logger) Info(template string, args ...interface{}) {
	pterm.ThemeDefault.InfoMessageStyle.Printfln(template, args...)
}

func (l *Logger) Error(template string, args ...interface{}) {
	pterm.ThemeDefault.ErrorMessageStyle.Printfln(template, args...)
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
