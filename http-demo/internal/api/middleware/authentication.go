package middleware

import (
	"net/http"
	"strings"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/labstack/echo/v4"
)

// AuthenticationMiddleware is the http middleware to handle authentication.
type AuthenticationMiddleware struct {
	config *config.Config
}

// NewAuthenticationMiddleware returns a new [AuthenticationMiddleware].
func NewAuthenticationMiddleware(config *config.Config) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		config: config,
	}
}

// Handle handles the http request authentication.
func (m *AuthenticationMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			logger := log.CtxLogger(request.Context())

			if m.config.GetBool("config.authentication.enabled") {
				if strings.TrimPrefix(request.Header.Get("authorization"), "Bearer ") != m.config.GetString("config.authentication.secret") {
					logger.Warn().Msg("authentication failed")

					return echo.NewHTTPError(http.StatusUnauthorized, "authentication failed")
				}

				logger.Info().Msg("authentication success")
			}

			return next(c)
		}
	}
}
