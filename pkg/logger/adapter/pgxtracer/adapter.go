package pgxtracer

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	maxDuration time.Duration
}

func NewAdapter(maxDuration time.Duration) *Logger {
	return &Logger{
		maxDuration: maxDuration,
	}
}

func (l *Logger) Query(ctx context.Context, sql string, duration time.Duration, rowsAffected int64, err error) {
	entry := zerolog.
		Ctx(ctx)

	var event *zerolog.Event

	switch {
	case err != nil:
		event = entry.Error().Err(err)
	case duration > l.maxDuration:
		event = entry.Warn()
	default:
		event = entry.Debug()
	}

	event = event.Str("sql", sql)

	if rowsAffected > 0 {
		event = event.Int64("rows", rowsAffected)
	}

	event.Msg("SQL")
}
