package tail

import "log"

type LogFunc func(string, ...interface{})

type Logger struct {
	Debugf    LogFunc
	Infof     LogFunc
	Noticef   LogFunc
	Warningf  LogFunc
	Errorf    LogFunc
	Criticalf LogFunc
}

func NullLogger(string, ...interface{}) {
}

func NewLogger() *Logger {
	return &Logger{
		Debugf:    NullLogger,
		Infof:     NullLogger,
		Noticef:   NullLogger,
		Warningf:  log.Printf,
		Errorf:    log.Printf,
		Criticalf: log.Fatalf,
	}
}

func NewLoggerLevel(level int) *Logger {
	l := NewLogger()

	if level >= 1 {
		l.Noticef = log.Printf
	}
	if level >= 2 {
		l.Infof = log.Printf
	}
	if level >= 4 {
		l.Debugf = log.Printf
	}
	return l
}
