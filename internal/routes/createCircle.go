package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/util"
)

func CreateCircleGet(c echo.Context) error {
    cc := c.(*util.CirclesContext)

	session, _ := cc.Store.Get(cc.Request(), "auth")

	if session.IsNew {
		return cc.Redirect(http.StatusSeeOther, "/signin")
	}

	_, ok := session.Values["first"].(string)
    _, ok = session.Values["email"].(string)

	if !ok {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

    return cc.Render(http.StatusOK, "createCircle.html", map[string]interface{}{
        "Route": "/create-circle",
    })
}

func CreateCirclePost(c echo.Context) error {
    cc := c.(*util.CirclesContext)

    session, _ := cc.Store.Get(cc.Request(), "auth")

    if session.IsNew {
        cc.Response().Header().Add("Hx-Redirect", "/signin")
		return cc.String(http.StatusOK, "")
    }

    email, ok := session.Values["email"].(string)
    if !ok {
        cc.Response().Header().Add("Hx-Redirect", "/signin")
		return cc.String(http.StatusOK, "")
    }

    name := cc.FormValue("name")

    fmt.Println(name)

    err := cc.DB.InsertCircle(name, email)
    if err != nil {
        return err
    }

    cc.Response().Header().Add("Hx-Redirect", "/home")
    return cc.String(http.StatusOK, "")
}
