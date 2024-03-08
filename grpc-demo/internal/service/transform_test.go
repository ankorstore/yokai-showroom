package service_test

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-showroom/grpc-demo/internal"
	"github.com/ankorstore/yokai-showroom/grpc-demo/internal/service"
	"github.com/ankorstore/yokai-showroom/grpc-demo/proto"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestTransformText(t *testing.T) {
	service.TransformerCounter.Reset()

	var grpcServer *grpc.Server
	var lis *bufconn.Listener
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	internal.RunTest(t, fx.Populate(&grpcServer, &lis, &logBuffer, &traceExporter, &metricsRegistry))

	defer func() {
		err := lis.Close()
		assert.NoError(t, err)

		grpcServer.GracefulStop()
	}()

	// client preparation
	conn, err := prepareGrpcClientTestConnection(lis)
	assert.NoError(t, err)

	client := proto.NewTransformTextServiceClient(conn)

	// call
	response, err := client.TransformText(context.Background(), &proto.TransformTextRequest{
		Transformer: proto.Transformer_TRANSFORMER_UPPERCASE,
		Text:        "test",
	})
	assert.NoError(t, err)

	// response assertions
	assert.Equal(t, "TEST", response.Text)

	// log assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "TransformText: test -> TEST",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(t, traceExporter, "TransformText")

	// metrics assertion
	expectedMetric := `
		# HELP transformer_total Total of TransformTextService transformer usage
		# TYPE transformer_total counter
		transformer_total{transformer="uppercase"} 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"transformer_total",
	)
	assert.NoError(t, err)
}

func TestTransformAndSplitText(t *testing.T) {
	service.TransformerCounter.Reset()

	var grpcServer *grpc.Server
	var lis *bufconn.Listener
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	internal.RunTest(t, fx.Populate(&grpcServer, &lis, &logBuffer, &traceExporter, &metricsRegistry))

	defer func() {
		err := lis.Close()
		assert.NoError(t, err)

		grpcServer.GracefulStop()
	}()

	// client preparation
	conn, err := prepareGrpcClientTestConnection(lis)
	assert.NoError(t, err)

	client := proto.NewTransformTextServiceClient(conn)

	// call
	stream, err := client.TransformAndSplitText(context.Background())
	assert.NoError(t, err)

	wait := make(chan struct{})

	go func() {
		err = stream.Send(&proto.TransformTextRequest{
			Transformer: proto.Transformer_TRANSFORMER_LOWERCASE,
			Text:        "THIS IS A TEST",
		})
		assert.NoError(t, err)

		err = stream.CloseSend()
		assert.NoError(t, err)
	}()

	var responses []*proto.TransformTextResponse
	go func() {
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			assert.NoError(t, err)

			responses = append(responses, resp)
		}

		close(wait)
	}()

	<-wait

	assert.Len(t, responses, 4)
	assert.Equal(t, "this", responses[0].Text)
	assert.Equal(t, "is", responses[1].Text)
	assert.Equal(t, "a", responses[2].Text)
	assert.Equal(t, "test", responses[3].Text)

	// logs assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "TransformTextAndSplit: -> THIS IS A TEST",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "TransformTextAndSplit: <- this",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "TransformTextAndSplit: <- is",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "TransformTextAndSplit: <- a",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "TransformTextAndSplit: <- test",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(t, traceExporter, "TransformAndSplitText")

	span, err := traceExporter.Span("TransformAndSplitText")
	assert.NoError(t, err)

	assert.Equal(t, "send word: this", span.Events[0].Name)
	assert.Equal(t, "send word: is", span.Events[1].Name)
	assert.Equal(t, "send word: a", span.Events[2].Name)
	assert.Equal(t, "send word: test", span.Events[3].Name)

	// metrics assertion
	expectedMetric := `
		# HELP transformer_total Total of TransformTextService transformer usage
		# TYPE transformer_total counter
		transformer_total{transformer="lowercase"} 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"transformer_total",
	)
	assert.NoError(t, err)
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
