package server

import (
	"io"

	"github.com/ankorstore/yokai-showroom/http-demo/module/fxmysqlmemory/config"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// ServerOptions are options for the [MySQLMemoryServerFactory] implementations.
type ServerOptions struct {
	Config    *config.MySQLMemoryServerConfig
	LogOutput io.Writer
	Tracer    trace.Tracer
}

// DefaultMySQLMemoryServerOptions are the default options used in the [DefaultMySQLMemoryServerFactory].
func DefaultMySQLMemoryServerOptions() ServerOptions {
	return ServerOptions{
		Config:    config.NewMySQLMemoryServerConfig(),
		LogOutput: io.Discard,
		Tracer:    noop.NewTracerProvider().Tracer(""),
	}
}

// MySQLMemoryServerOption are functional options for the [MySQLMemoryServerFactory] implementations.
type MySQLMemoryServerOption func(o *ServerOptions)

func WithConfig(config *config.MySQLMemoryServerConfig) MySQLMemoryServerOption {
	return func(o *ServerOptions) {
		if config != nil {
			o.Config = config
		}
	}
}

func WithLogOutput(output io.Writer) MySQLMemoryServerOption {
	return func(o *ServerOptions) {
		if output != nil {
			o.LogOutput = output
		}
	}
}

func WithTracer(tracer trace.Tracer) MySQLMemoryServerOption {
	return func(o *ServerOptions) {
		if tracer != nil {
			o.Tracer = tracer
		}
	}
}
