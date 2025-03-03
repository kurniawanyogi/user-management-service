package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"user-management-service/common"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const SkipFrame = 1

// StandardLogger is standard logger struct type
type StandardLogger struct {
	ZeroLogger zerolog.Logger
}

type Config struct {
	AppName string `json:"appName" yaml:"appName"`
	Debug   bool   `json:"debug" yaml:"debug"`
}

// Logger is logger singleton
var Logger StandardLogger

// Init to initiate Logger
func Init(config Config) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	Logger.ZeroLogger = zerolog.New(zerolog.SyncWriter(os.Stdout)).With().
		Caller().
		Timestamp().
		Str("app_name", config.AppName).
		Logger()
}

func writeZeroLog(ev *zerolog.Event, tags ...Tag) *zerolog.Event {
	for _, t := range tags {
		ev = ev.Str(t.Key, fmt.Sprintf("%+v", t.Value))
	}
	return ev
}

func getStackTrace(from, to int) (traces []string) {
	for from < to {
		if _, f, l, ok := runtime.Caller(from); ok {
			traces = append(traces, fmt.Sprintf("called from %s:%d", f, l))
		}
		from++
	}
	return traces
}

// Debug to provide global zerolog log Debug
func Debug(ctx context.Context, msg string, tags ...Tag) {
	writeZeroLog(
		Logger.ZeroLogger.Debug().CallerSkipFrame(SkipFrame),
		append(tags, GetAllLoggingTagInTagStr(ctx)...)...).
		Msg(msg)
}

// Info to provide global zerolog log Info
func Info(ctx context.Context, msg string, tags ...Tag) {
	writeZeroLog(
		Logger.ZeroLogger.Info().CallerSkipFrame(SkipFrame),
		append(tags, GetAllLoggingTagInTagStr(ctx)...)...).
		Msg(msg)
}

// Warn to provide global zerolog log Warn
func Warn(ctx context.Context, msg string, tags ...Tag) {
	writeZeroLog(
		Logger.ZeroLogger.Warn().CallerSkipFrame(SkipFrame),
		append(tags, GetAllLoggingTagInTagStr(ctx)...)...,
	).Msg(msg)
}

// Error to provide global zerolog log Error
func Error(ctx context.Context, msg string, err error, tags ...Tag) {
	ev := writeZeroLog(
		Logger.ZeroLogger.Error().CallerSkipFrame(SkipFrame),
		append(tags, GetAllLoggingTagInTagStr(ctx)...)...,
	)
	if err != nil {
		ev.Err(err)
	}
	ev.Strs("stack_trace", getStackTrace(2, 5)).Msg(msg)
}

// Fatal to provide global zerolog log Fatal
func Fatal(ctx context.Context, msg string, tags ...Tag) {
	writeZeroLog(
		Logger.ZeroLogger.Fatal().CallerSkipFrame(SkipFrame),
		append(tags, GetAllLoggingTagInTagStr(ctx)...)...,
	).Msg(msg)
}

func NewContextFromParent(c *gin.Context) context.Context {
	requestID := c.Request.Context().Value(common.XRequestIdHeader).(string)

	return AddRequestID(
		context.WithValue(
			context.Background(),
			common.XRequestIdHeader,
			requestID,
		),
		requestID,
	)
}
