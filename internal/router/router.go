package router

import (
	"github.com/Satishcg12/multicommers/internal/router/routes"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to Echomers")
	})

	// group routes
	api := e.Group("/api")
	{
		routes.RegisterVendorAuthRoutes(api)

	}

}
