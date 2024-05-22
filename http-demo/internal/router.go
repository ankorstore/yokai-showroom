package internal

import (
	"github.com/ankorstore/yokai-showroom/http-demo/internal/handler"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/handler/gopher"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/middleware"
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

// Router is used to register the application HTTP middlewares and handlers.
func Router() fx.Option {
	return fx.Options(
		// authentication middleware
		fxhttpserver.AsMiddleware(middleware.NewAuthenticationMiddleware, fxhttpserver.GlobalUse),
		// dashboard handler
		fxhttpserver.AsHandler("GET", "", handler.NewDashboardHandler),
		// gophers CRUD handlers group
		fxhttpserver.AsHandlersGroup(
			"/gophers",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("GET", "", gopher.NewListGophersHandler),
				fxhttpserver.NewHandlerRegistration("POST", "", gopher.NewCreateGopherHandler),
				fxhttpserver.NewHandlerRegistration("GET", "/:id", gopher.NewGetGopherHandler),
				fxhttpserver.NewHandlerRegistration("DELETE", "/:id", gopher.NewDeleteGopherHandler),
			},
		),
	)
}
