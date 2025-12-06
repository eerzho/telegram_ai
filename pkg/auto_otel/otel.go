package autootel

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type OTel struct {
	tracerProvider *trace.TracerProvider
	loggerProvider *log.LoggerProvider
	meterProvider  *metric.MeterProvider
}

func Setup(ctx context.Context) (*OTel, error) {
	res, err := newResource(ctx)
	if err != nil {
		return nil, err
	}

	traceExporter, err := newTraceExporter(ctx)
	if err != nil {
		return nil, err
	}
	tracerProvider := newTracerProvider(res, traceExporter)
	otel.SetTracerProvider(tracerProvider)

	logExporter, err := newLogExporter(ctx)
	if err != nil {
		return nil, err
	}
	loggerProvider := newLoggerProvider(res, logExporter)
	global.SetLoggerProvider(loggerProvider)

	metricExporter, err := newMetricExporter(ctx)
	if err != nil {
		return nil, err
	}
	meterProvider := newMeterProvider(res, metricExporter)
	otel.SetMeterProvider(meterProvider)

	if err := runtimeStart(); err != nil {
		return nil, err
	}

	if err := hostStart(); err != nil {
		return nil, err
	}

	return &OTel{
		tracerProvider: tracerProvider,
		loggerProvider: loggerProvider,
		meterProvider:  meterProvider,
	}, nil
}

func (o *OTel) Shutdown(ctx context.Context) error {
	errs := make([]error, 0, 3)
	if err := o.tracerProvider.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}
	if err := o.loggerProvider.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}
	if err := o.meterProvider.Shutdown(ctx); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
