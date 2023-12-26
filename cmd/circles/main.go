package circles

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vincer2040/circles/internal/render"
)

func Main() {

	e := echo.New()

	e.Renderer = render.New()

	e.Use(middleware.Logger())

	e.GET("/", rootGet)

	e.Logger.Fatal(e.Start(":6969"))
}

func rootGet(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
