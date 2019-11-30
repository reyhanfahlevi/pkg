package logger

import (
	"fmt"
)

// IsValid check if Logger is created using constructor
func (l *Logger) IsValid() bool {
	return l.valid
}

// Debug function
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug().Timestamp().Msg(fmt.Sprint(args...))
}

// Debugln function
func (l *Logger) Debugln(args ...interface{}) {
	l.logger.Debug().Timestamp().Msg(fmt.Sprintln(args...))
}

// Debugf function
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Timestamp().Msgf(format, v...)
}

// DebugWithFields function
func (l *Logger) DebugWithFields(msg string, kv map[string]interface{}) {
	l.logger.Debug().Timestamp().Fields(kv).Msg(msg)
}

// Info function
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info().Timestamp().Msg(fmt.Sprint(args...))
}

// Infoln function
func (l *Logger) Infoln(args ...interface{}) {
	l.logger.Info().Timestamp().Msg(fmt.Sprintln(args...))
}

// Infof function
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info().Timestamp().Msgf(format, v...)
}

// InfoWithFields function
func (l *Logger) InfoWithFields(msg string, kv map[string]interface{}) {
	l.logger.Info().Timestamp().Fields(kv).Msg(msg)
}

// Warn function
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn().Timestamp().Msg(fmt.Sprint(args...))
}

// Warnln function
func (l *Logger) Warnln(args ...interface{}) {
	l.logger.Warn().Timestamp().Msg(fmt.Sprintln(args...))
}

// Warnf function
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Timestamp().Msgf(format, v...)
}

// WarnWithFields function
func (l *Logger) WarnWithFields(msg string, kv map[string]interface{}) {
	l.logger.Warn().Timestamp().Fields(kv).Msg(msg)
}

// Error function
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error().Timestamp().Msg(fmt.Sprint(args...))
}

// Errorln function
func (l *Logger) Errorln(args ...interface{}) {
	l.logger.Error().Timestamp().Msg(fmt.Sprintln(args...))
}

// Errorf function
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Timestamp().Msgf(format, v...)
}

// ErrorWithFields function
func (l *Logger) ErrorWithFields(msg string, kv map[string]interface{}) {
	l.logger.Error().Timestamp().Fields(kv).Msg(msg)
}

// Errors function to log errors package
func (l *Logger) Errors(err error) {
	l.logger.Error().Timestamp().Msg(err.Error())
}

// Fatal function
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal().Timestamp().Msg(fmt.Sprint(args...))
}

// Fatalln function
func (l *Logger) Fatalln(args ...interface{}) {
	l.logger.Fatal().Timestamp().Msg(fmt.Sprintln(args...))
}

// Fatalf function
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Timestamp().Msgf(format, v...)
}

// FatalWithFields function
func (l *Logger) FatalWithFields(msg string, kv map[string]interface{}) {
	l.logger.Fatal().Timestamp().Fields(kv).Msg(msg)
}
