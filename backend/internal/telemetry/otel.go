package telemetry

import (
    "context"

    "github.com/example/team-stack/backend/internal/config"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTracer(cfg *config.Config) (*sdktrace.TracerProvider, func(context.Context) error) {
    exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(cfg.AppName),
        )),
    )
    otel.SetTracerProvider(tp)
    return tp, tp.Shutdown
}
