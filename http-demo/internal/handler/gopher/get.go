package gopher

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/service"
	"github.com/labstack/echo/v4"
)

// GetGopherHandler is the http handler to get a gopher.
type GetGopherHandler struct {
	service *service.GopherService
}

// NewGetGopherHandler returns a new GetGopherHandler.
func NewGetGopherHandler(service *service.GopherService) *GetGopherHandler {
	return &GetGopherHandler{
		service: service,
	}
}

// Handle handles the http request.
func (h *GetGopherHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid gopher id: %v", err))
		}

		gopher, err := h.service.Get(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("cannot get gopher with id %d: %v", id, err))
			}

			return fmt.Errorf("cannot get gopher: %w", err)
		}

		return c.JSON(http.StatusOK, gopher)
	}
}
