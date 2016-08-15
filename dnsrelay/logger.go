package dnsrelay

import (
	"github.com/apsdehal/go-logger"
	"os"
	"io"
	"fmt"
	"strings"
)



// Level defines all available log levels for log messages.
type LogLevel int

// Log levels.
const (
	CRITICAL LogLevel = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var levelNames = []string{
	"CRITICAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
}

// String returns the string representation of a logging level.
func (p LogLevel) String() string {
	return levelNames[p]
}

func (self *LogLevel) UnmarshalTOML(data []byte) (err error) {
	level := string(data)
	level = strings.TrimSpace(level)
	level = strings.Trim(level, "\"")

	for i, name := range levelNames {
		if name == level {
			*self = LogLevel(i)
			return
		}
	}
	return fmt.Errorf("logger: invalid log level:%s\nvalid choice:%v\n", level, levelNames)
}

type Logger struct {
	log   *logger.Logger
	level LogLevel
}

func NewLogger(path string, name string, level LogLevel) (log *Logger, err error) {
	var f io.WriteCloser
	var color int

	if path == "" {
		f = os.Stdout
		color = 1
	} else {
		color = 0
		f, err = os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			return
		}
	}

	l, err := logger.New(name, color, f)
	if err != nil {
		return
	}

	return &Logger{log:l, level:level}, nil
}


// IsEnabledFor will return true if logging is enabled for the given module.
func (l *Logger) IsEnabledFor(level LogLevel) bool {
	return level <= l.level
}

func (l *Logger) WriteLog(level LogLevel, msg string) {
	if l.IsEnabledFor(level) {
		l.log.Log(level.String(), msg)
	}
}

func (l *Logger) WriteLogf(level LogLevel, f string, args... interface{}) {
	if l.IsEnabledFor(level) {
		s := fmt.Sprintf(f, args...)
		l.log.Log(level.String(), s)
	}
}


// Panic is just like func l.Critical except that it is followed by a call to panic
func (l *Logger) Panic(message string) {
	l.WriteLog(CRITICAL, message)
	panic(message)
}

// Info logs a message at Info level
func (l *Logger) Panicf(f string, args... interface{}) {
	l.WriteLogf(CRITICAL, f, args...)
	panic("")
}

// Error logs a message at Error level
func (l *Logger) Error(message string) {
	l.WriteLog(ERROR, message)
}

// Info logs a message at Info level
func (l *Logger) Errorf(f string, args... interface{}) {
	l.WriteLogf(ERROR, f, args...)
}


// Warning logs a message at Warning level
func (l *Logger) Warning(message string) {
	l.WriteLog(WARNING, message)
}

// Info logs a message at Info level
func (l *Logger) Warningf(f string, args... interface{}) {
	l.WriteLogf(WARNING, f, args...)
}

// Info logs a message at Info level
func (l *Logger) Info(message string) {
	l.WriteLog(INFO, message)
}


// Info logs a message at Info level
func (l *Logger) Infof(f string, args... interface{}) {
	l.WriteLogf(INFO, f, args...)

}

// Debug logs a message at Debug level
func (l *Logger) Debug(message string) {
	l.WriteLog(DEBUG, message)
}

// Debug logs a message at Debug level
func (l *Logger) Debugf(f string, args... interface{}) {
	l.WriteLogf(DEBUG, f, args...)

}