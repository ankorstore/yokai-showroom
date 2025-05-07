package gopher

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/domain"
	"github.com/labstack/echo/v4"
)

// GetGopherHandler is the http handler to get a gopher.
type GetGopherHandler struct {
	service *domain.GopherService
}

// NewGetGopherHandler returns a new GetGopherHandler.
func NewGetGopherHandler(service *domain.GopherService) *GetGopherHandler {
	return &GetGopherHandler{
		service: service,
	}
}

// Handle handles the http request.
func (h *GetGopherHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		gopherId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid gopher id: %v", err))
		}

		gopher, err := h.service.Get(c.Request().Context(), gopherId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("cannot find gopher with id %d: %v", gopherId, err))
			}

			return err
		}

		return c.JSON(http.StatusOK, gopher)
	}
}
