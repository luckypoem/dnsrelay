package dnsrelay

import (
	"github.com/apsdehal/go-logger"
	"os"
	"io"
	"fmt"
)

type Logger struct {
	log *logger.Logger
}

func NewLogger(path string, name string) (log *Logger, err error) {
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

	return &Logger{log:l}, nil
}


// Fatal is just like func l,Cr.tical logger except that it is followed by exit to program
func (l *Logger) Fatal(message string) {
	l.log.Log("CRITICAL", message)
	os.Exit(1)
}

// Info logs a message at Info level
func (l *Logger) Fatalf(f string, args... interface{}) {
	s := fmt.Sprintf(f, args...)
	l.log.Log("CRITICAL", s)
	os.Exit(1)
}

// Panic is just like func l.Critical except that it is followed by a call to panic
func (l *Logger) Panic(message string) {
	l.log.Log("CRITICAL", message)
	panic(message)
}

// Info logs a message at Info level
func (l *Logger) Panicf(f string, args... interface{}) {
	s := fmt.Sprintf(f, args...)
	l.log.Log("CRITICAL", s)
	panic(s)
}

// Critical logs a message at a Critical Level
func (l *Logger) Critical(message string) {
	l.log.Log("CRITICAL", message)
}

// Info logs a message at Info level
func (l *Logger) Criticalf(f string, args... interface{}) {
	s := fmt.Sprintf(f, args...)
	l.log.Log("CRITICAL", s)
}


// Error logs a message at Error level
func (l *Logger) Error(message string) {
	l.log.Log("ERROR", message)
}

// Info logs a message at Info level
func (l *Logger) Errorf(f string, args... interface{}) {
	s := fmt.Sprintf(f, args...)
	l.log.Log("ERROR", s)
}


// Warning logs a message at Warning level
func (l *Logger) Warning(message string) {
	l.log.Log("WARNING", message)
}

// Info logs a message at Info level
func (l *Logger) Warningf(f string, args... interface{}) {
	s := fmt.Sprintf(f, args...)
	l.log.Log("WARNING", s)
}


// Notice logs a message at Notice level
func (l *Logger) Notice(message string) {
	l.log.Log("NOTICE", message)
}

// Info logs a message at Info level
func (l *Logger) Noticef(f string, args... interface{}) {
	s := fmt.Sprintf(f, args...)
	l.log.Log("NOTICE", s)
}


// Info logs a message at Info level
func (l *Logger) Info(message string) {
	l.log.Log("INFO", message)
}


// Info logs a message at Info level
func (l *Logger) Infof(f string, args... interface{}) {
	s := fmt.Sprintf(f, args...)
	l.log.Log("INFO", s)
}

// Debug logs a message at Debug level
func (l *Logger) Debug(message string) {
	l.log.Log("DEBUG", message)
}

// Debug logs a message at Debug level
func (l *Logger) Debugf(f string, args... interface{}) {
	s := fmt.Sprintf(f, args...)
	l.log.Log("DEBUG", s)
}