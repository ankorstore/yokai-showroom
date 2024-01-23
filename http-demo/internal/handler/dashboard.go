package handler

import (
	"net/http"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

// DashboardHistogram is an example of histogram metric for the dashboard rendering.
var DashboardHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "dashboard_duration_seconds",
	Help:    "The duration of the dashboard rendering in seconds",
	Buckets: prometheus.DefBuckets,
})

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
		start := time.Now()
		defer func() {
			DashboardHistogram.Observe(time.Since(start).Seconds())
		}()

		return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
			"title": h.config.GetString("config.dashboard.title"),
		})
	}
}
