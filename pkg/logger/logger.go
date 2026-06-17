package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(skipFrameCount).Logger()

	return &Logger{
		logger: &logger,
	}
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.log(l.logger.Debug(), message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(l.logger.Info(), message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(l.logger.Warn(), message, args...)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.log(l.logger.Error(), message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.log(l.logger.Fatal(), message, args...)
	os.Exit(1)
}

func (l *Logger) log(e *zerolog.Event, message interface{}, args ...interface{}) {
	if len(args) == 0 {
		switch msg := message.(type) {
		case error:
			e.Msg(msg.Error())
		case string:
			e.Msg(msg)
		default:
			e.Msgf("%v", message)
		}
	} else {
		switch msg := message.(type) {
		case string:
			e.Msgf(msg, args...)
		default:
			e.Msgf(fmt.Sprintf("%v", message), args...)
		}
	}
}
