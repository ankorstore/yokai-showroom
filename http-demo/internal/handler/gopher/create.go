package gopher

import (
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/service"
	"github.com/labstack/echo/v4"
)

// CreateGopherHandler is the http handler to create a gopher.
type CreateGopherHandler struct {
	service *service.GopherService
}

// NewCreateGopherHandler returns a new CreateGopherHandler.
func NewCreateGopherHandler(service *service.GopherService) *CreateGopherHandler {
	return &CreateGopherHandler{
		service: service,
	}
}

// Handle handles the http request.
func (h *CreateGopherHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		gopher := new(model.Gopher)
		if err := c.Bind(gopher); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("cannot bind gopher: %v", err))
		}

		err := h.service.Create(c.Request().Context(), gopher)
		if err != nil {
			return fmt.Errorf("cannot create gopher: %w", err)
		}

		return c.JSON(http.StatusCreated, gopher)
	}
}
