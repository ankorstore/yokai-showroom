package internal

import (
	"github.com/ankorstore/yokai-showroom/grpc-demo/internal/interceptor"
	"github.com/ankorstore/yokai-showroom/grpc-demo/internal/service"
	"github.com/ankorstore/yokai-showroom/grpc-demo/proto"
	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/fxmetrics"
	"go.uber.org/fx"
)

// ProvideServices is used to register the application services.
func ProvideServices() fx.Option {
	return fx.Options(
		// gRPC server interceptors
		fxgrpcserver.AsGrpcServerUnaryInterceptor(interceptor.NewAuthenticationUnaryInterceptor),
		fxgrpcserver.AsGrpcServerStreamInterceptor(interceptor.NewAuthenticationStreamInterceptor),
		// gRPC server service
		fxgrpcserver.AsGrpcServerService(service.NewTransformTextService, &proto.TransformTextService_ServiceDesc),
		// metrics
		fxmetrics.AsMetricsCollector(service.TransformerCounter),
	)
}
