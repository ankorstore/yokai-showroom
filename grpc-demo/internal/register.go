package internal

import (
	"github.com/ankorstore/yokai-showroom/grpc-demo/internal/interceptor"
	"github.com/ankorstore/yokai-showroom/grpc-demo/internal/service"
	"github.com/ankorstore/yokai-showroom/grpc-demo/proto"
	"github.com/ankorstore/yokai/fxgrpcserver"
	"github.com/ankorstore/yokai/fxmetrics"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// interceptors
		fxgrpcserver.AsGrpcServerUnaryInterceptor(interceptor.NewAuthenticationUnaryInterceptor),
		fxgrpcserver.AsGrpcServerStreamInterceptor(interceptor.NewAuthenticationStreamInterceptor),
		// service
		fxgrpcserver.AsGrpcServerService(service.NewTransformTextService, &proto.TransformTextService_ServiceDesc),
		// metrics
		fxmetrics.AsMetricsCollector(service.TransformerCounter),
	)
}
