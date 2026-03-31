package logging

import (
	"testing"
)

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level LogLevel
		want  string
	}{
		{All, "all"},
		{Debug, "debug"},
		{Info, "info"},
		{Error, "error"},
		{LogLevel(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("LogLevel.String() = %v, want %v", got, tt.want)
		}
	}
}

func TestLogLevelFromString(t *testing.T) {
	tests := []struct {
		input string
		want  LogLevel
	}{
		{"all", All},
		{"debug", Debug},
		{"info", Info},
		{"error", Error},
		{"unknown", All},
	}

	for _, tt := range tests {
		if got := LogLevelFromString(tt.input); got != tt.want {
			t.Errorf("LogLevelFromString(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestNewLogger(t *testing.T) {
	l := NewLogger()
	if l == nil {
		t.Fatal("NewLogger() returned nil")
	}
	if l.LogLevel != Info {
		t.Errorf("expected default LogLevel Info, got %v", l.LogLevel)
	}
	if l.InfoWriter == nil || l.ErrWriter == nil || l.DebugWriter == nil {
		t.Error("writers should not be nil")
	}
}

func TestLoggerMethods(t *testing.T) {
	l := NewLogger()
	
	// Just verify they don't crash
	l.LogLevel = All
	l.Debug("test debug %s", "arg")
	l.Info("test info %s", "arg")
	l.Error("test error %s", "arg")

	// Verify writers
	_, _ = l.InfoWriter.Write([]byte("info test"))
	_, _ = l.ErrWriter.Write([]byte("error test"))
	_, _ = l.DebugWriter.Write([]byte("debug test"))
}
