package zerogoose

import "github.com/rs/zerolog"

type Logger struct {
	logger *zerolog.Logger
}

func NewLogger(l *zerolog.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Msgf(format, v...)
}

func (l Logger) Printf(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}
