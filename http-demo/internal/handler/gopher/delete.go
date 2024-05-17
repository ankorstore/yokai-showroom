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

// DeleteGopherHandler is the http handler to delete a gopher.
type DeleteGopherHandler struct {
	service *service.GopherService
}

// NewDeleteGopherHandler returns a new DeleteGopherHandler.
func NewDeleteGopherHandler(service *service.GopherService) *DeleteGopherHandler {
	return &DeleteGopherHandler{
		service: service,
	}
}

// Handle handles the http request.
func (h *DeleteGopherHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid gopher id: %v", err))
		}

		err = h.service.Delete(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("cannot get gopher with id %d: %v", id, err))
			}

			return fmt.Errorf("cannot delete gopher: %w", err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
