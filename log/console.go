package log

import (
	"log"
	"os"
	"runtime"
)

type Brush func(string) string

func init() {
	Register("console", NewConsole)
}

func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []Brush{
	NewBrush("1;36"), // Trace      cyan
	NewBrush("1;34"), // Debug      blue
	NewBrush("1;32"), // Info       green
	NewBrush("1;33"), // Warn       yellow
	NewBrush("1;31"), // Error      red
	NewBrush("1;35"), // Critical   purple
	NewBrush("1;31"), // Fatal      red
}

// ConsoleWriter implements interface LogProvider and writes messages to terminal.
type ConsoleWriter struct {
	Log     *log.Logger
}

// create ConsoleWriter returning as LoggerInterface.
func NewConsole() LogProvider {
	return &ConsoleWriter{
		Log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (cw *ConsoleWriter) Write(l *Logger, msg string, level int) error {
	if l.LogLevel > level {
		return nil
	}
	if runtime.GOOS == "windows" {
		cw.Log.Println(msg)
	} else {
		cw.Log.Println(colors[level](msg))
	}
	return nil
}
