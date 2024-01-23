package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/worker-demo/internal"
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
)

type PublisherMock struct {
	mock.Mock
}

func (m *PublisherMock) Publish(context.Context, string) error {
	return m.Called().Error(0)
}

func TestPublishHandlerSuccess(t *testing.T) {
	publisherMock := new(PublisherMock)
	publisherMock.On("Publish").Return(nil)

	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(
		t,
		fx.Populate(
			&httpServer,
			&logBuffer,
			&traceExporter,
		),
		fx.Replace(
			fx.Annotate(
				publisherMock,
				fx.As(new(service.Publisher)),
			),
		),
	)

	// call [GET] /publish?message=test
	req := httptest.NewRequest(http.MethodGet, "/publish?message=test", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	// response assertion
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "publication of message test success.")

	// logs assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "worker-demo-app",
		"method":  "GET",
		"status":  200,
		"uri":     "/publish?message=test",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /publish",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/publish"),
		semconv.HTTPStatusCode(http.StatusOK),
	)
}

func TestPublishHandlerError(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	publisherMock := new(PublisherMock)
	publisherMock.On("Publish").Return(fmt.Errorf("custom error"))

	internal.RunTest(
		t,
		fx.Populate(
			&httpServer,
			&logBuffer,
			&traceExporter,
		),
		fx.Replace(
			fx.Annotate(
				publisherMock,
				fx.As(new(service.Publisher)),
			),
		),
	)

	// call [GET] /publish?message=test
	req := httptest.NewRequest(http.MethodGet, "/publish?message=test", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	// response assertion
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "publication of message test failure: custom error.")

	// logs assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "worker-demo-app",
		"method":  "GET",
		"status":  500,
		"uri":     "/publish?message=test",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"GET /publish",
		semconv.HTTPMethod(http.MethodGet),
		semconv.HTTPRoute("/publish"),
		semconv.HTTPStatusCode(http.StatusInternalServerError),
	)
}
