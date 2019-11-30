package log

import (
	"github.com/reyhanfahlevi/pkg/go/log/logger"
)

// Level logger
type Level logger.Level

// Logger interface
type Logger interface {
	SetLevel(level logger.Level)
	Debug(args ...interface{})
	Debugln(args ...interface{})
	Debugf(format string, args ...interface{})
	DebugWithFields(msg string, kv map[string]interface{})
	Info(args ...interface{})
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})
	InfoWithFields(msg string, kv map[string]interface{})
	Warn(args ...interface{})
	Warnln(args ...interface{})
	Warnf(format string, args ...interface{})
	WarnWithFields(msg string, kv map[string]interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorWithFields(msg string, kv map[string]interface{})
	Errors(err error)
	Fatal(args ...interface{})
	Fatalln(args ...interface{})
	Fatalf(format string, args ...interface{})
	FatalWithFields(msg string, kv map[string]interface{})
	IsValid() bool // IsValid check if Logger is created using constructor
}

// Level option
const (
	DebugLevel = Level(logger.DebugLevel)
	InfoLevel  = Level(logger.InfoLevel)
	WarnLevel  = Level(logger.WarnLevel)
	ErrorLevel = Level(logger.ErrorLevel)
	FatalLevel = Level(logger.FatalLevel)
)

// Debug prints debug level log like log.Print
func Debug(args ...interface{}) {
	debugLogger.Debug(args...)
}

// Debugln prints debug level log like log.Println
func Debugln(args ...interface{}) {
	debugLogger.Debugln(args...)
}

// Debugf prints debug level log like log.Printf
func Debugf(format string, v ...interface{}) {
	debugLogger.Debugf(format, v...)
}

// DebugWithFields prints debug level log with additional fields.
// useful when output is in json format
func DebugWithFields(msg string, fields map[string]interface{}) {
	debugLogger.DebugWithFields(msg, fields)
}

// Print info level log like log.Print
func Print(v ...interface{}) {
	infoLogger.Info(v...)
}

// Println info level log like log.Println
func Println(v ...interface{}) {
	infoLogger.Infoln(v...)
}

// Printf info level log like log.Printf
func Printf(format string, v ...interface{}) {
	infoLogger.Infof(format, v...)
}

// Info prints info level log like log.Print
func Info(args ...interface{}) {
	infoLogger.Info(args...)
}

// Infoln prints info level log like log.Println
func Infoln(args ...interface{}) {
	infoLogger.Infoln(args...)
}

// Infof prints info level log like log.Printf
func Infof(format string, v ...interface{}) {
	infoLogger.Infof(format, v...)
}

// InfoWithFields prints info level log with additional fields.
// useful when output is in json format
func InfoWithFields(msg string, fields map[string]interface{}) {
	infoLogger.InfoWithFields(msg, fields)
}

// Warn prints warn level log like log.Print
func Warn(args ...interface{}) {
	warnLogger.Warn(args...)
}

// Warnln prints warn level log like log.Println
func Warnln(args ...interface{}) {
	warnLogger.Warnln(args...)
}

// Warnf prints warn level log like log.Printf
func Warnf(format string, v ...interface{}) {
	warnLogger.Warnf(format, v...)
}

// WarnWithFields prints warn level log with additional fields.
// useful when output is in json format
func WarnWithFields(msg string, fields map[string]interface{}) {
	warnLogger.WarnWithFields(msg, fields)
}

// Error prints error level log like log.Print
func Error(args ...interface{}) {
	errLogger.Error(args...)
}

// Errorln prints error level log like log.Println
func Errorln(args ...interface{}) {
	errLogger.Errorln(args...)
}

// Errorf prints error level log like log.Printf
func Errorf(format string, v ...interface{}) {
	errLogger.Errorf(format, v...)
}

// ErrorWithFields prints error level log with additional fields.
// useful when output is in json format
func ErrorWithFields(msg string, fields map[string]interface{}) {
	errLogger.ErrorWithFields(msg, fields)
}

// Errors can handle error from tdk/x/go/errors package
func Errors(err error) {
	errLogger.Errors(err)
}

// Fatal prints fatal level log like log.Print
func Fatal(args ...interface{}) {
	fatalLogger.Fatal(args...)
}

// Fatalln prints fatal level log like log.Println
func Fatalln(args ...interface{}) {
	fatalLogger.Fatalln(args...)
}

// Fatalf prints fatal level log like log.Printf
func Fatalf(format string, v ...interface{}) {
	fatalLogger.Fatalf(format, v...)
}

// FatalWithFields prints fatal level log with additional fields.
// useful when output is in json format
func FatalWithFields(msg string, fields map[string]interface{}) {
	fatalLogger.FatalWithFields(msg, fields)
}
