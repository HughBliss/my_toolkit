package tracer

import (
	"fmt"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

type hook struct{}

func HookForLogger() zerolog.Hook {
	return new(hook)
}

func (z *hook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	ctx := e.GetCtx()
	if ctx == nil {
		return
	}
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	ev := fmt.Sprintf("%s", reflect.ValueOf(e).Elem().FieldByName("buf"))

	span.AddEvent(level.String(), trace.WithAttributes(
		attribute.String("meta", ev+"}"),
		attribute.String("message", message),
	))
	if level == zerolog.ErrorLevel {
		span.SetStatus(codes.Error, message)
	}

	if span.SpanContext().HasTraceID() {
		e = e.Str("traceID", span.SpanContext().TraceID().String())
	}

	if span.SpanContext().HasSpanID() {
		e = e.Str("spanID", span.SpanContext().SpanID().String())
	}
}
