package internal

import (
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/handler"
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

func ProvideRouting() fx.Option {
	return fx.Options(
		fxhttpserver.AsHandler("GET", "/publish", handler.NewPublishHandler),
	)
}
