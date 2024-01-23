package gopher

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// UpdateGopherHandler is the http handler to update a gopher.
type UpdateGopherHandler struct {
	service *service.GopherService
}

// NewUpdateGopherHandler returns a new UpdateGopherHandler.
func NewUpdateGopherHandler(service *service.GopherService) *UpdateGopherHandler {
	return &UpdateGopherHandler{
		service: service,
	}
}

// Handle handles the http request.
func (h *UpdateGopherHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid gopher id: %v", err))
		}

		update := new(model.Gopher)
		if err = c.Bind(update); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("cannot bind gopher: %v", err))
		}

		gopher, err := h.service.Update(c.Request().Context(), id, update)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("cannot get gopher with id %d: %v", id, err))
			}

			return fmt.Errorf("cannot update gopher: %w", err)
		}

		return c.JSON(http.StatusOK, gopher)
	}
}
