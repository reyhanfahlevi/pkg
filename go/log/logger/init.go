package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// list of log level
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Log level
const (
	DebugLevelString = "debug"
	InfoLevelString  = "info"
	WarnLevelString  = "warn"
	ErrorLevelString = "error"
	FatalLevelString = "fatal"
)

const DefaultTimeFormat = time.RFC3339

// Level of log
type Level int

type Logger struct {
	logger zerolog.Logger
	config Config
	valid  bool
}

type Config struct {
	Level      Level
	AppName    string
	LogFile    string
	TimeFormat string
	CallerSkip int
	Caller     bool
	UseColor   bool
	UseJSON    bool
}

func New(config *Config) (*Logger, error) {
	if config == nil {
		config = &Config{
			Level:      InfoLevel,
			TimeFormat: DefaultTimeFormat,
		}
	}

	if config.TimeFormat == "" {
		config.TimeFormat = DefaultTimeFormat
	}

	lgr, err := newLogger(*config)
	if err != nil {
		return nil, err
	}
	l := Logger{
		logger: lgr,
		config: *config,
		valid:  true,
	}
	return &l, nil
}

// SetLevel for setting log level
func (l *Logger) SetLevel(level Level) {
	if level < DebugLevel || level > FatalLevel {
		level = InfoLevel
	}
	if level != l.config.Level {
		l.logger = setLevel(l.logger, level)
		l.config.Level = level
	}
}

func newLogger(config Config) (zerolog.Logger, error) {
	var (
		lgr zerolog.Logger
	)

	zerolog.TimeFieldFormat = config.TimeFormat
	zerolog.CallerSkipFrameCount = 4 + config.CallerSkip

	var writers zerolog.LevelWriter
	if config.UseJSON {
		writers = zerolog.MultiLevelWriter(os.Stderr)
	} else {
		writers = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    !config.UseColor,
			TimeFormat: config.TimeFormat,
		})
	}

	file, err := config.OpenLogFile()
	if err != nil {
		return lgr, err
	} else if file != nil {
		writers = zerolog.MultiLevelWriter(writers, file)
	}

	if config.AppName != "" {
		lgr = zerolog.New(writers).With().Str("appname", config.AppName).Logger()
	}

	lgr = setLevel(lgr, config.Level)
	if config.Caller {
		lgr = lgr.With().Caller().Logger()
	}

	return lgr, nil
}

func setLevel(lgr zerolog.Logger, level Level) zerolog.Logger {
	switch level {
	case DebugLevel:
		lgr = lgr.Level(zerolog.DebugLevel)
	case InfoLevel:
		lgr = lgr.Level(zerolog.InfoLevel)
	case WarnLevel:
		lgr = lgr.Level(zerolog.WarnLevel)
	case ErrorLevel:
		lgr = lgr.Level(zerolog.ErrorLevel)
	case FatalLevel:
		lgr = lgr.Level(zerolog.FatalLevel)
	default:
		lgr = lgr.Level(zerolog.InfoLevel)
	}
	return lgr
}

// StringToLevel to set string to level
func StringToLevel(level string) Level {
	switch strings.ToLower(level) {
	case DebugLevelString:
		return DebugLevel
	case InfoLevelString:
		return InfoLevel
	case WarnLevelString:
		return WarnLevel
	case ErrorLevelString:
		return ErrorLevel
	case FatalLevelString:
		return FatalLevel
	default:
		return InfoLevel
	}
}

// OpenLogFile tries to open the log file (creates it if not exists) in write-only/append mode and return it
// Note: the func return nil for both *os.File and error if the file name is empty string
func (c *Config) OpenLogFile() (*os.File, error) {
	if c.LogFile == "" {
		return nil, nil
	}

	err := os.MkdirAll(filepath.Dir(c.LogFile), 0755)
	if err != nil && err != os.ErrExist {
		return nil, err
	}

	return os.OpenFile(c.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
}
