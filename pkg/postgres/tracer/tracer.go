package tracer

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type (
	tracerStartQueryContextKey = struct{}
	tracerSQLQueryContextKey   = struct{}
	tracerArgsQueryContextKey  = struct{}
)

type Logger interface {
	Query(ctx context.Context, sql string, duration time.Duration, rowsAffected int64, err error)
}

type LogTracer struct {
	logger Logger
}

func NewLogTracer(logger Logger) *LogTracer {
	return &LogTracer{logger: logger}
}

func (t *LogTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, tracerStartQueryContextKey{}, time.Now())
	ctx = context.WithValue(ctx, tracerSQLQueryContextKey{}, data.SQL)

	return context.WithValue(ctx, tracerArgsQueryContextKey{}, data.Args)
}

func (t *LogTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	start, ok := ctx.Value(tracerStartQueryContextKey{}).(time.Time)
	if !ok {
		return
	}
	sql, ok := ctx.Value(tracerSQLQueryContextKey{}).(string)
	if !ok {
		return
	}
	args, ok := ctx.Value(tracerArgsQueryContextKey{}).([]any)
	if !ok {
		return
	}

	interpolatedSQL := t.cleanQuery(t.inlineQueryWithArgs(sql, args))
	duration := time.Since(start)
	rowsAffected := data.CommandTag.RowsAffected()

	t.logger.Query(ctx, interpolatedSQL, duration, rowsAffected, data.Err)
}

func (t *LogTracer) cleanQuery(input string) string {
	cleaner := func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' {
			return -1
		}

		return r
	}

	return strings.Map(cleaner, input)
}

func (t *LogTracer) inlineQueryWithArgs(sql string, args []any) string {
	if len(args) == 0 {
		return sql
	}

	const nullValue = "null"
	for i, arg := range args {
		argType := reflect.TypeOf(arg)
		argVal := reflect.ValueOf(arg)
		var value string
		switch {
		case argType == nil:
			value = nullValue

		case argType == reflect.TypeOf(time.Time{}):
			timeArg := arg.(time.Time)
			value = fmt.Sprintf("'%s'", timeArg.Format("2006-01-02 15:04:05"))

		case argType.Kind() == reflect.Ptr:
			if arg == nil || argVal.IsNil() {
				value = nullValue
			} else {
				elemVal := argVal.Elem()
				switch elemVal.Kind() {
				case reflect.String:
					value = fmt.Sprintf("'%v'", elemVal.Interface())
				case reflect.Struct:
					if elemVal.Type() == reflect.TypeFor[time.Time]() {
						value = fmt.Sprintf("'%s'", elemVal.Interface().(time.Time).Format("2006-01-02 15:04:05"))
					} else {
						value = fmt.Sprintf("'%v'", elemVal.Interface())
					}
				default:
					value = fmt.Sprintf("%v", elemVal.Interface())
				}
			}
		case !argType.ConvertibleTo(reflect.TypeFor[int64]()) &&
			!argType.ConvertibleTo(reflect.TypeFor[float64]()) && argType.Kind() != reflect.Bool:
			value = fmt.Sprintf("'%v'", arg)
		default:
			value = fmt.Sprintf("%v", arg)
		}

		placeholder := fmt.Sprintf("$%d", i+1)
		sql = strings.Replace(sql, placeholder, value, 1)
	}

	return sql
}
