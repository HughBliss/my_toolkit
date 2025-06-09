package reporter

import (
	"context"
	zfg "github.com/chaindead/zerocfg"
	"github.com/hughbliss/my_toolkit/tracer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.elastic.co/ecszerolog"
	"go.opentelemetry.io/otel/trace"
	"os"
)

var (
	LogLevel = zfg.Str("log_level", "DEBUG", "LOGLEVEL уровень логирования")
)

func Init(hooks ...zerolog.Hook) {
	level, _ := zerolog.ParseLevel(*LogLevel)
	zerolog.SetGlobalLevel(level)
	log.Logger = ecszerolog.New(os.Stdout).Hook(hooks...)
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
		l.Debug().Msg("end")
		span.End()
	}
}

func InitReporter(serviceName string, hooks ...zerolog.Hook) Reporter {
	l := log.Logger.With().Str("service", serviceName).Logger().Hook(hooks...)
	l.Debug().Msg("init")
	return &reporter{
		l: l,
		t: tracer.Provider.Tracer(serviceName),
	}
}
