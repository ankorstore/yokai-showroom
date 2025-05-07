package interceptor_test

import (
	"context"
	"io"
	"testing"

	"github.com/ankorstore/yokai-showroom/grpc-demo/internal"
	"github.com/ankorstore/yokai-showroom/grpc-demo/proto"
	"github.com/ankorstore/yokai/grpcserver/grpcservertest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//nolint:containedctx
func TestAuthenticationUnaryInterceptor(t *testing.T) {
	t.Setenv("AUTH_ENABLED", "true")
	t.Setenv("AUTH_SECRET", "valid-secret")

	var connFactory grpcservertest.TestBufconnConnectionFactory

	internal.RunTest(t, fx.Populate(&connFactory))

	// conn preparation
	conn, err := connFactory.Create(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	defer func() {
		err = conn.Close()
		assert.NoError(t, err)
	}()

	// client preparation
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
				Text: "test",
			})
			assert.Equal(t, tData.want, got)
		})
	}
}

//nolint:containedctx
func TestAuthenticationStreamInterceptor(t *testing.T) {
	t.Setenv("AUTH_ENABLED", "true")
	t.Setenv("AUTH_SECRET", "valid-secret")

	var connFactory grpcservertest.TestBufconnConnectionFactory

	internal.RunTest(t, fx.Populate(&connFactory))

	// conn preparation
	conn, err := connFactory.Create(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	defer func() {
		err = conn.Close()
		assert.NoError(t, err)
	}()

	// client preparation
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
