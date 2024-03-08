package interceptor_test

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/ankorstore/yokai-showroom/grpc-demo/internal"
	"github.com/ankorstore/yokai-showroom/grpc-demo/proto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestAuthenticationUnaryInterceptor(t *testing.T) {
	t.Setenv("AUTH_ENABLED", "true")
	t.Setenv("AUTH_SECRET", "valid-secret")

	var grpcServer *grpc.Server
	var lis *bufconn.Listener

	internal.RunTest(t, fx.Populate(&grpcServer, &lis))

	defer func() {
		err := lis.Close()
		assert.NoError(t, err)

		grpcServer.GracefulStop()
	}()

	// client preparation
	conn, err := prepareGrpcClientTestConnection(lis)
	assert.NoError(t, err)

	client := proto.NewTransformTextServiceClient(conn)

	// tests cases
	testCases := map[string]struct {
		ctx  context.Context
		want error
	}{
		"auth success": {
			metadata.AppendToOutgoingContext(context.Background(), "authorization", "valid-secret"),
			nil,
		},
		"auth failure with missing authorization": {
			context.Background(),
			status.Errorf(codes.InvalidArgument, "missing authorization"),
		},
		"auth failure with invalid token": {
			metadata.AppendToOutgoingContext(context.Background(), "authorization", "invalid-secret"),
			status.Errorf(codes.Unauthenticated, "invalid token"),
		},
	}

	// tests assertions
	for tName, tData := range testCases {
		t.Run(tName, func(t *testing.T) {
			_, got := client.TransformText(tData.ctx, &proto.TransformTextRequest{
				Transformer: proto.Transformer_TRANSFORMER_UPPERCASE,
				Text:        "test",
			})
			assert.Equal(t, tData.want, got)
		})
	}
}

func TestAuthenticationStreamInterceptor(t *testing.T) {
	t.Setenv("AUTH_ENABLED", "true")
	t.Setenv("AUTH_SECRET", "valid-secret")

	var grpcServer *grpc.Server
	var lis *bufconn.Listener

	internal.RunTest(t, fx.Populate(&grpcServer, &lis))

	defer func() {
		err := lis.Close()
		assert.NoError(t, err)

		grpcServer.GracefulStop()
	}()

	// client preparation
	conn, err := prepareGrpcClientTestConnection(lis)
	assert.NoError(t, err)

	client := proto.NewTransformTextServiceClient(conn)

	// tests cases
	testCases := map[string]struct {
		ctx           context.Context
		wantAtConnect error
		wantAtReceive error
	}{
		"auth success": {
			metadata.AppendToOutgoingContext(context.Background(), "authorization", "valid-secret"),
			nil,
			io.EOF,
		},
		"auth failure with missing authorization": {
			context.Background(),
			nil,
			status.Errorf(codes.InvalidArgument, "missing authorization"),
		},
		"auth failure with invalid token": {
			metadata.AppendToOutgoingContext(context.Background(), "authorization", "invalid-secret"),
			nil,
			status.Errorf(codes.Unauthenticated, "invalid token"),
		},
	}

	// tests assertions
	for tName, tData := range testCases {
		t.Run(tName, func(t *testing.T) {
			stream, gotAtConnect := client.TransformAndSplitText(tData.ctx)
			assert.Equal(t, tData.wantAtConnect, gotAtConnect)

			if stream != nil {
				err = stream.CloseSend()
				assert.NoError(t, err)
			}

			_, gotAtReceive := stream.Recv()
			assert.Equal(t, tData.wantAtReceive, gotAtReceive)
		})
	}
}

func prepareGrpcClientTestConnection(lis *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		"",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
