package log

import (
	"errors"

	"github.com/reyhanfahlevi/pkg/go/log/logger"
)

// Config of log
type Config struct {
	// Level log level, default: debug
	Level Level

	// Format of the log's time field
	// default RFC3339="2006-01-02T15:04:05Z07:00"
	TimeFormat string

	// AppName is the application name this log belong to
	AppName string

	// Caller, option to print caller line numbers.
	// make sure you understand the overhead when use this
	Caller bool

	// LogFile is output file for log other than debug log
	// this is not needed by default,
	// application is expected to run in containerized environment
	LogFile string

	// DebugFile is output file for debug log
	// this is not needed by default,
	// application is expected to run in containerized environment
	DebugFile string

	// UseColor, option to colorize log in console.
	UseColor bool

	// UseJSON, option to print in json format.
	UseJSON bool
}

var (
	infoLogger, _  = NewLogger(&Config{Level: Level(logger.InfoLevel), UseColor: false})
	debugLogger, _ = NewLogger(&Config{Level: Level(logger.DebugLevel), UseColor: false})
	warnLogger     = infoLogger
	errLogger      = infoLogger
	fatalLogger    = infoLogger
	loggers        = [5]*Logger{
		&debugLogger,
		&infoLogger,
		&warnLogger,
		&errLogger,
		&fatalLogger,
	}
)

// NewLogger will create new logger with specific configuration
// the return must be set using set logger to register it
func NewLogger(config *Config) (Logger, error) {
	l, err := logger.New(&logger.Config{
		Level:      logger.Level(config.Level),
		AppName:    config.AppName,
		LogFile:    config.LogFile,
		TimeFormat: config.TimeFormat,
		Caller:     config.Caller,
		UseColor:   config.UseColor,
		UseJSON:    config.UseJSON,
	})
	if err != nil {
		return nil, err
	}

	return l, nil
}

// SetLogger will set the logger configuration on certain level
func SetLogger(level Level, lgr Logger) error {
	if level < DebugLevel || level > FatalLevel {
		return errors.New("invalid level")
	}
	if lgr == nil || !lgr.IsValid() {
		return errors.New("invalid logger")
	}
	*loggers[level] = lgr
	return nil
}

// SetLevel adjusts log level threshold.
// Only log with level higher or equal with this level will be printed
func SetLevel(level Level) {
	if level < 0 {
		level = InfoLevel
	}

	debugLogger.SetLevel(logger.Level(level))
	infoLogger.SetLevel(logger.Level(level))
	warnLogger.SetLevel(logger.Level(level))
	errLogger.SetLevel(logger.Level(level))
	fatalLogger.SetLevel(logger.Level(level))
}

// SetConfig creates new default (info & debug) logger based on given config
func SetConfig(config *Config) error {
	var (
		newDebugLogger    Logger
		newLogger         Logger
		err               error
		debugLoggerConfig = logger.Config{Level: logger.DebugLevel}
		loggerConfig      = logger.Config{Level: logger.InfoLevel}
	)

	if config != nil {
		loggerConfig = logger.Config{
			Level:      logger.Level(config.Level),
			AppName:    config.AppName,
			LogFile:    config.LogFile,
			TimeFormat: config.TimeFormat,
			Caller:     config.Caller,
			UseColor:   config.UseColor,
			UseJSON:    config.UseJSON,
		}

		debugLoggerConfig = loggerConfig
		debugLoggerConfig.LogFile = config.DebugFile
	}

	newLogger, err = logger.New(&loggerConfig)
	if err != nil {
		return err
	}

	if newLogger != nil {
		infoLogger = newLogger
		warnLogger = newLogger
		errLogger = newLogger
		fatalLogger = newLogger
	}

	newDebugLogger, err = logger.New(&debugLoggerConfig)
	if err != nil {
		return err
	}
	if newDebugLogger != nil {
		debugLogger = newDebugLogger
	}

	return nil
}
