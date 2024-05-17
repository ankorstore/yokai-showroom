package gopher

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/service"
	"github.com/labstack/echo/v4"
)

type CreateGopherParams struct {
	Name string `json:"name" form:"name" query:"name"`
	Job  string `json:"job" form:"job" query:"job"`
}

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
		ctx := c.Request().Context()

		params := new(CreateGopherParams)
		if err := c.Bind(params); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid parameters: %v", err))
		}

		gopherId, err := h.service.Create(ctx, params.Name, params.Job)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("cannot create gopher: %v", err))
		}

		gopher, err := h.service.Get(ctx, gopherId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("owner with id %d not found", gopherId))
			}

			return err
		}

		return c.JSON(http.StatusCreated, gopher)
	}
}
