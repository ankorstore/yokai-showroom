package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
	"github.com/labstack/echo/v4"
)

type PublishHandler struct {
	publisher service.Publisher
}

func NewPublishHandler(publisher service.Publisher) *PublishHandler {
	return &PublishHandler{
		publisher: publisher,
	}
}

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
