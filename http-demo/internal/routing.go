package internal

import (
	"github.com/ankorstore/yokai-showroom/http-demo/internal/handler"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/handler/gopher"
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

func ProvideRouting() fx.Option {
	return fx.Options(
		fxhttpserver.AsHandler("GET", "", handler.NewDashboardHandler),
		// gophers CRUD
		fxhttpserver.AsHandlersGroup(
			"/gophers",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "", gopher.NewListGophersHandler),
				fxhttpserver.NewHandlerRegistration("POST", "", gopher.NewCreateGopherHandler),
				fxhttpserver.NewHandlerRegistration("GET", "/:id", gopher.NewGetGopherHandler),
				fxhttpserver.NewHandlerRegistration("PATCH", "/:id", gopher.NewUpdateGopherHandler),
				fxhttpserver.NewHandlerRegistration("DELETE", "/:id", gopher.NewDeleteGopherHandler),
			},
		),
	)
}
