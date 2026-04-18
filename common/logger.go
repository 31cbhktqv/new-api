package common

import (
	"fmt"
	"io"
	"os"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[LogLevel]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError:  "ERROR",
}

type Logger struct {
	writer    io.Writer
	minLevel  LogLevel
}

// Changed default min level to Debug so I can see all log output during local development
var DefaultLogger = NewLogger(os.Stdout, LevelDebug)

func NewLogger(w io.Writer, minLevel LogLevel) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{writer: w, minLevel: minLevel}
}

func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.minLevel {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	name := levelNames[level]
	formatted := fmt.Sprintf(msg, args...)
	fmt.Fprintf(l.writer, "[%s] [%s] %s\n", timestamp, name, formatted)
}

func (l *Logger) Debug(msg string, args ...interface{}) { l.log(LevelDebug, msg, args...) }
func (l *Logger) Info(msg string, args ...interface{})  { l.log(LevelInfo, msg, args...) }
func (l *Logger) Warn(msg string, args ...interface{})  { l.log(LevelWarn, msg, args...) }
func (l *Logger) Error(msg string, args ...interface{}) { l.log(LevelError, msg, args...) }

func Debug(msg string, args ...interface{}) { DefaultLogger.Debug(msg, args...) }
func Info(msg string, args ...interface{})  { DefaultLogger.Info(msg, args...) }
func Warn(msg string, args ...interface{})  { DefaultLogger.Warn(msg, args...) }
func Error(msg string, args ...interface{}) { DefaultLogger.Error(msg, args...) }
