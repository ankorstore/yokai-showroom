package gopher

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/service"
	"github.com/labstack/echo/v4"
)

// ListGophersHandler is the http handler to list all gophers.
type ListGophersHandler struct {
	service *service.GopherService
}

// NewListGophersHandler returns a new ListGophersHandler.
func NewListGophersHandler(service *service.GopherService) *ListGophersHandler {
	return &ListGophersHandler{
		service: service,
	}
}

// Handle handles the http request.
func (h *ListGophersHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		gophers, err := h.service.List(c.Request().Context())
		if err != nil {
			return fmt.Errorf("cannot list gophers: %w", err)
		}

		return c.JSON(http.StatusOK, gophers)
	}
}
