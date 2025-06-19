package trace_middleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

func AddTraceIDToResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		spanCtx := trace.SpanContextFromContext(ctx)
		if !spanCtx.IsValid() {
			return next(c)
		}
		if spanCtx.HasTraceID() {
			c.Response().Header().Add("x-trace-id", spanCtx.TraceID().String())
		}
		return next(c)
	}
}
