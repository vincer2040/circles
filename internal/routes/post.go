package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/util"
)

func PostPost(c echo.Context) error {
    cc := c.(*util.CirclesContext)

    session, _ := cc.Store.Get(c.Request(), "auth")
    circle := cc.Param("circle")

    author, ok := session.Values["email"].(string)
    if !ok {
        return cc.Redirect(http.StatusSeeOther, "/signin")
    }

    description := cc.FormValue("description")

    err := cc.DB.InsertPost(circle, author, description)
    if err != nil {
        return err
    }

    return cc.String(http.StatusOK, "")
}
