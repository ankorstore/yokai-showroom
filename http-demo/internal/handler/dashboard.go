package handler

import (
	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

// DashboardHandler is the http handler to render the dashboard.
type DashboardHandler struct {
	config *config.Config
}

// NewDashboardHandler returns a new [DashboardHandler].
func NewDashboardHandler(config *config.Config) *DashboardHandler {
	return &DashboardHandler{
		config: config,
	}
}

// Handle handles the http request.
func (h *DashboardHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
			"title": h.config.GetString("config.dashboard.title"),
		})
	}
}
