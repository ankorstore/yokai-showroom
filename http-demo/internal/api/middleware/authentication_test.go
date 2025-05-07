package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestAuthenticationDisabled(t *testing.T) {
	t.Setenv("AUTH_ENABLED", "false")

	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer))

	// [GET] /
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	// response assertion
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthenticationSuccess(t *testing.T) {
	t.Setenv("AUTH_ENABLED", "true")
	t.Setenv("AUTH_SECRET", "valid-secret")

	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer))

	// [GET] /
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-secret")
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	// response assertion
	assert.Equal(t, http.StatusOK, rec.Code)

	// log assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "authentication success",
	})
}

func TestAuthenticationFailure(t *testing.T) {
	t.Setenv("AUTH_ENABLED", "true")
	t.Setenv("AUTH_SECRET", "valid-secret")

	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer))

	// [GET] /
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-secret")
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	// response assertion
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// log assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "warn",
		"message": "authentication failed",
	})
}
