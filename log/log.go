package log

import (
	"fmt"
	"os"
)

var (
	Log *Logger
)

func init() {
	Log = &Logger{
		Prefix:  "nanogit",
		Adapter: "console",
	}
}

func Trace(format string, v ...interface{}) {
	Log.Trace(format, v...)
}

func Debug(format string, v ...interface{}) {
	Log.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	Log.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	Log.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	Log.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	Error(format, v...)
	os.Exit(1)
}

// ———————————————————————————————————————————————————————————————————
//                      Log interface

const (
	TRACE = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
	FATAL
)

type LogProvider interface {
	Write(l *Logger, msg string, level int) error
}

var adapters = make(map[string]func() LogProvider)

// Registers given logger provider to adapters.
func Register(name string, log func() LogProvider) {
	if log == nil {
		panic("log: register provider is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("log: register called twice for provider \"" + name + "\"")
	}
	adapters[name] = log
}

type logMsg struct {
	Level   int
	Message string
}

// Default logger. It can contain several providers and log message into all providers.
type Logger struct {
	Prefix   string
	LogLevel int
	Adapter  string
}

func (l *Logger) writerMsg(level int, msg string) {
	provider, has := adapters[l.Adapter]
	if !has {
		panic("log: no registered adapter: " + l.Adapter)
	}
	provider().Write(Log, msg, level)
}

func (l *Logger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf("trace: "+format, v...)
	l.writerMsg(TRACE, msg)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("debug: "+format, v...)
	l.writerMsg(DEBUG, msg)
}

func (l *Logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("info: "+format, v...)
	l.writerMsg(INFO, msg)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("warning: "+format, v...)
	l.writerMsg(WARN, msg)
}

func (l *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("error: "+format, v...)
	l.writerMsg(ERROR, msg)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf("fatal: "+format, v...)
	l.writerMsg(FATAL, msg)
	os.Exit(1)
}
