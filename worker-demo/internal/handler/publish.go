package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// PublishHandler is the http handler for [GET] /publish.
type PublishHandler struct {
	messageChannel chan string
}

type PublishHandlerParam struct {
	fx.In
	MessageChannel chan string `name:"pub-sub-message-channel"`
}

// NewPublishHandler returns a new PublishHandler.
func NewPublishHandler(p PublishHandlerParam) *PublishHandler {
	return &PublishHandler{
		messageChannel: p.MessageChannel,
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

		h.messageChannel <- message

		return c.String(http.StatusAccepted, fmt.Sprintf("publication of message %s success.", message))
	}
}
