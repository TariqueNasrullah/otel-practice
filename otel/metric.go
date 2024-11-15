package otel

import (
	"context"
	"errors"
	"github.com/TariqueNasrullah/otel-practice/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func Init(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(ctx context.Context) error

	shutdown = func(ctx context.Context) error {
		var shutdownErr error
		for _, fn := range shutdownFuncs {
			shutdownErr = errors.Join(shutdownErr, fn(ctx))
		}
		shutdownFuncs = nil
		return shutdownErr
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	re, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL, semconv.ServiceName(config.AppName),
	))
	if err != nil {
		handleErr(err)
		return
	}

	meterProvider, err := initMetric(re)
	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func initMetric(resource *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := prometheus.New(prometheus.WithNamespace(config.AppName))
	if err != nil {
		return nil, err
	}

	return metric.NewMeterProvider(
		metric.WithReader(metricExporter),
		metric.WithResource(resource),
	), nil
}
