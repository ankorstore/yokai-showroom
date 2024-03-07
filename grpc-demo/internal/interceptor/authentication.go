package interceptor

import (
	"context"
	"strings"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthenticationUnaryInterceptor is an authentication unary interceptor.
type AuthenticationUnaryInterceptor struct {
	config *config.Config
}

// NewAuthenticationUnaryInterceptor returns a new [AuthenticationUnaryInterceptor].
func NewAuthenticationUnaryInterceptor(cfg *config.Config) *AuthenticationUnaryInterceptor {
	return &AuthenticationUnaryInterceptor{
		config: cfg,
	}
}

// HandleUnary handles the unary request authentication.
func (i *AuthenticationUnaryInterceptor) HandleUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger := log.CtxLogger(ctx)

		if i.config.GetBool("config.authentication.enabled") {
			err := authenticate(ctx, i.config.GetString("config.authentication.secret"))
			if err != nil {
				logger.Warn().Msg("unary authentication failed")

				return nil, err
			}

			logger.Info().Msg("unary authentication success")
		}

		return handler(ctx, req)
	}
}

// AuthenticationStreamInterceptor is an authentication unary interceptor.
type AuthenticationStreamInterceptor struct {
	config *config.Config
}

// NewAuthenticationStreamInterceptor returns a new [AuthenticationStreamInterceptor].
func NewAuthenticationStreamInterceptor(cfg *config.Config) *AuthenticationStreamInterceptor {
	return &AuthenticationStreamInterceptor{
		config: cfg,
	}
}

// HandleStream handles the stream request authentication.
func (i *AuthenticationStreamInterceptor) HandleStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger := log.CtxLogger(ss.Context())

		if i.config.GetBool("config.authentication.enabled") {
			err := authenticate(ss.Context(), i.config.GetString("config.authentication.secret"))
			if err != nil {
				logger.Warn().Msg("stream authentication failed")

				return err
			}

			logger.Info().Msg("stream authentication success")
		}

		return handler(srv, ss)
	}
}

func authenticate(ctx context.Context, expectedToken string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	authMd, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing authorization")
	}

	if len(authMd) < 1 {
		return status.Errorf(codes.InvalidArgument, "missing token")
	}

	if strings.TrimPrefix(authMd[0], "Bearer ") != expectedToken {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return nil
}
