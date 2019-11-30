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

// Debug function
func Debug(args ...interface{}) {
	debugLogger.Debug(args...)
}

// Debugln function
func Debugln(args ...interface{}) {
	debugLogger.Debugln(args...)
}

// Debugf function
func Debugf(format string, v ...interface{}) {
	debugLogger.Debugf(format, v...)
}

// DebugWithFields function
func DebugWithFields(msg string, kv map[string]interface{}) {
	debugLogger.DebugWithFields(msg, kv)
}

// Info function
func Info(args ...interface{}) {
	infoLogger.Info(args...)
}

// Infoln function
func Infoln(args ...interface{}) {
	infoLogger.Infoln(args...)
}

// Infof function
func Infof(format string, v ...interface{}) {
	infoLogger.Infof(format, v...)
}

// InfoWithFields function
func InfoWithFields(msg string, kv map[string]interface{}) {
	infoLogger.InfoWithFields(msg, kv)
}

// Warn function
func Warn(args ...interface{}) {
	warnLogger.Warn(args...)
}

// Warnln function
func Warnln(args ...interface{}) {
	warnLogger.Warnln(args...)
}

// Warnf function
func Warnf(format string, v ...interface{}) {
	warnLogger.Warnf(format, v...)
}

// WarnWithFields function
func WarnWithFields(msg string, kv map[string]interface{}) {
	warnLogger.WarnWithFields(msg, kv)
}

// Error function
func Error(args ...interface{}) {
	errLogger.Error(args...)
}

// Errorln function
func Errorln(args ...interface{}) {
	errLogger.Errorln(args...)
}

// Errorf function
func Errorf(format string, v ...interface{}) {
	errLogger.Errorf(format, v...)
}

// ErrorWithFields function
func ErrorWithFields(msg string, kv map[string]interface{}) {
	errLogger.ErrorWithFields(msg, kv)
}

// Errors function to log errors package
func Errors(err error) {
	errLogger.Errors(err)
}

// Fatal function
func Fatal(args ...interface{}) {
	fatalLogger.Fatal(args...)
}

// Fatalln function
func Fatalln(args ...interface{}) {
	fatalLogger.Fatalln(args...)
}

// Fatalf function
func Fatalf(format string, v ...interface{}) {
	fatalLogger.Fatalf(format, v...)
}

// FatalWithFields function
func FatalWithFields(msg string, kv map[string]interface{}) {
	fatalLogger.FatalWithFields(msg, kv)
}
