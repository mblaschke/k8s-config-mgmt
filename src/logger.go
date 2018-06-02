package main

import (
	"log"
	"os"
	"fmt"
)

const (
	LogPrefix         = ""
	prefixMain        = ":: "
	prefixCategory    = "   "
	prefixSubCategory = "     "
	prefixStep        = "     - "
	prefixStepResult  = "        -> "
	prefixErr  = "[ERROR] "
)

type DaemonLogger struct {
	*log.Logger
}

var (
	Verbose bool
)

func CreateDaemonLogger(flags int) *DaemonLogger {
	return &DaemonLogger{log.New(os.Stdout, LogPrefix, flags)}
}

func CreateDaemonErrorLogger(flags int) *DaemonLogger {
	return &DaemonLogger{log.New(os.Stderr, LogPrefix, flags)}
}

func (l *DaemonLogger) Verbose(message string, sprintf ...interface{}) {
	if Verbose {
		if len(sprintf) > 0 {
			message = fmt.Sprintf(message, sprintf...)
		}

		l.Println(message)
	}
}

func (l *DaemonLogger) Main(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	l.Println(prefixMain + message)
}


func (l *DaemonLogger) Category(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	l.Println(prefixCategory + message)
}

func (l *DaemonLogger) SubCategory(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	l.Println(prefixSubCategory + message)
}

func (l *DaemonLogger) Step(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	l.Println(prefixStep + message)
}

func (l *DaemonLogger) StepResult(message string, sprintf ...interface{}) {
	if len(sprintf) > 0 {
		message = fmt.Sprintf(message, sprintf...)
	}

	l.Println(prefixStepResult + message)
}

// Log error object as message
func (l *DaemonLogger) Error(msg string, err error) {
	l.Println(fmt.Sprintf("[ERROR] %v: %v", msg, err))
}
