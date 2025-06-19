package reporter

import (
	"context"
	zfg "github.com/chaindead/zerocfg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"os"
	"time"
)

var (
	group    = zfg.NewGroup("log")
	LogLevel = zfg.Str("level", "DEBUG", "LOG_LEVEL уровень логирования", zfg.Group(group))
)

func Init(appName string, appVer string, env string, hooks ...zerolog.Hook) {
	// Базовые настройки полей
	zerolog.ErrorFieldName = "error_message" // Убираем точки из названий полей
	zerolog.ErrorStackFieldName = "error_stack"
	zerolog.TimestampFieldName = "timestamp" // Стандартное поле для Loki
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	zerolog.CallerSkipFrameCount = 4
	level, _ := zerolog.ParseLevel(*LogLevel)
	zerolog.SetGlobalLevel(level)
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", appName). // Лейбл для Loki
		Str("environment", env). // Лейбл для Loki
		Str("version", appVer).  // Версия схемы логов
		Logger().Hook(hooks...)
}

type Reporter interface {
	Start(ctx context.Context, methodName string) (context.Context, zerolog.Logger, func())
}

type reporter struct {
	l zerolog.Logger
	t trace.Tracer
}

func (w reporter) Start(ctx context.Context, methodName string) (context.Context, zerolog.Logger, func()) {
	ctx, span := w.t.Start(ctx, methodName)
	l := w.l.With().Ctx(ctx).Str("method", methodName).Logger()
	l.Debug().Msg("start")

	return ctx, l, func() {
		if pnc := recover(); pnc != nil {
			if err, ok := pnc.(error); ok {
				l.Error().Err(err).Stack().Msg("panic")
			} else {
				l.Error().Any("payload", pnc).Msg("panic")
			}
		}
		l.Debug().Msg("end")
		span.End()
	}
}

func InitReporter(serviceName string, hooks ...zerolog.Hook) Reporter {
	l := log.Logger.With().Str("component", serviceName).Logger().Hook(hooks...)
	l.Debug().Msg("init")
	return &reporter{
		l: l,
		t: otel.GetTracerProvider().Tracer(serviceName),
	}
}
