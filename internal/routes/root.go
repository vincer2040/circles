package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/util"
)

func RootGet(c echo.Context) error {
    pageData := util.PageData {
        Route: "/",
    }
    return c.Render(http.StatusOK, "index.html", pageData)
}
