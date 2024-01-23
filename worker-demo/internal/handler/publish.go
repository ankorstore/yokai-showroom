package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
	"github.com/labstack/echo/v4"
)

// PublishHandler is the http handler for [GET] /publish.
type PublishHandler struct {
	publisher service.Publisher
}

// NewPublishHandler returns a new PublishHandler.
func NewPublishHandler(publisher service.Publisher) *PublishHandler {
	return &PublishHandler{
		publisher: publisher,
	}
}

// Handle handles the http request.
func (h *PublishHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		message := time.Now().Format(time.DateTime)

		messageParam := c.QueryParam("message")
		if messageParam != "" {
			message = messageParam
		}

		err := h.publisher.Publish(c.Request().Context(), message)
		if err != nil {
			return c.String(
				http.StatusInternalServerError,
				fmt.Sprintf("publication of message %s failure: %s.", message, err.Error()),
			)
		}

		return c.String(
			http.StatusOK,
			fmt.Sprintf("publication of message %s success.", message),
		)
	}
}
