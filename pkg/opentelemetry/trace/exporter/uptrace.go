package exporter

import (
	"context"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

type UpTraceConfig struct {
	Dsn string
}

func UpTraceExporter(config *UpTraceConfig) (tracesdk.SpanExporter, error) {
	dsn, err := uptrace.ParseDSN(config.Dsn)
	if err != nil {
		return nil, err
	}
	return otlptrace.New(context.Background(), otlpTraceClient(dsn))
}

func otlpTraceClient(dsn *uptrace.DSN) otlptrace.Client {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(dsn.OTLPHost()),
		otlptracegrpc.WithHeaders(map[string]string{
			// Set the Uptrace DSN here or use UPTRACE_DSN env var.
			"uptrace-dsn": dsn.String(),
		}),
		otlptracegrpc.WithCompressor(gzip.Name),
	}

	if dsn.Scheme == "https" {
		// Create credentials using system certificates.
		creds := credentials.NewClientTLSFromCert(nil, "")
		options = append(options, otlptracegrpc.WithTLSCredentials(creds))
	} else {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.NewClient(options...)
}
