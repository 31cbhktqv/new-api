package common

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger_InfoWritesOutput(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, LevelInfo)
	l.Info("hello %s", "world")
	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestLogger_DebugSuppressedAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, LevelInfo)
	l.Debug("this should not appear")
	if buf.Len() != 0 {
		t.Errorf("expected no output for debug at info level, got: %s", buf.String())
	}
}

func TestLogger_DebugVisibleAtDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, LevelDebug)
	l.Debug("debug message")
	if !strings.Contains(buf.String(), "DEBUG") {
		t.Errorf("expected DEBUG in output, got: %s", buf.String())
	}
}

func TestLogger_WarnAndError(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf, LevelDebug)
	l.Warn("warn msg")
	l.Error("error msg")
	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output")
	}
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output")
	}
}

func TestLogger_NilWriterDefaultsToStdout(t *testing.T) {
	l := NewLogger(nil, LevelInfo)
	if l.writer == nil {
		t.Error("expected non-nil writer when nil is passed")
	}
}

func TestDefaultLogger_PackageLevelFunctions(t *testing.T) {
	var buf bytes.Buffer
	DefaultLogger = NewLogger(&buf, LevelDebug)
	Info("package level info")
	if !strings.Contains(buf.String(), "INFO") {
		t.Errorf("expected INFO from package-level Info()")
	}
}
