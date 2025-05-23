package gopher

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/domain"
	"github.com/labstack/echo/v4"
)

// ListGophersHandler is the http handler to list all gophers.
type ListGophersHandler struct {
	service *domain.GopherService
}

// NewListGophersHandler returns a new ListGophersHandler.
func NewListGophersHandler(service *domain.GopherService) *ListGophersHandler {
	return &ListGophersHandler{
		service: service,
	}
}

// Handle handles the http request.
func (h *ListGophersHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		gophers, err := h.service.List(c.Request().Context(), c.QueryParam("name"), c.QueryParam("job"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot list gophers: %v", err))
		}

		return c.JSON(http.StatusOK, gophers)
	}
}
