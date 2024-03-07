package internal

import (
	"github.com/ankorstore/yokai-showroom/grpc-demo/internal/interceptor"
	"github.com/ankorstore/yokai-showroom/grpc-demo/internal/service"
	"github.com/ankorstore/yokai-showroom/grpc-demo/proto"
	"github.com/ankorstore/yokai/fxgrpcserver"
	"go.uber.org/fx"
)

// ProvideServices is used to register the application services.
func ProvideServices() fx.Option {
	return fx.Options(
		// interceptors
		fxgrpcserver.AsGrpcServerUnaryInterceptor(interceptor.NewAuthenticationUnaryInterceptor),
		fxgrpcserver.AsGrpcServerStreamInterceptor(interceptor.NewAuthenticationStreamInterceptor),
		// service
		fxgrpcserver.AsGrpcServerService(service.NewTransformTextServiceService, &proto.TransformTextService_ServiceDesc),
	)
}
