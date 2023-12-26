package circles

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vincer2040/circles/internal/render"
	"github.com/vincer2040/circles/internal/routes"
)

func Main() {

	e := echo.New()

	e.Renderer = render.New()

	e.Use(middleware.Logger())
	e.Static("/css", "public/styles")

	e.GET("/", routes.RootGet)

	e.Logger.Fatal(e.Start(":6969"))
}
