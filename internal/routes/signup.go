package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/user"
	"github.com/vincer2040/circles/internal/util"
)

func SignupGet(c echo.Context) error {
	pageData := util.PageData{
		Route: "/signup",
	}
	return c.Render(http.StatusOK, "signup.html", pageData)
}

func SignupPost(c echo.Context) error {
	cc := c.(*util.CirclesContext)
	session, _ := cc.Store.Get(c.Request(), "auth")
	if !session.IsNew {
		authed, ok := session.Values["authenticated"].(bool)
		if ok && authed {
			c.Redirect(http.StatusSeeOther, "/me")
		}
	}
	first := cc.FormValue("first")
	last := cc.FormValue("last")
	email := cc.FormValue("email")
	password := cc.FormValue("password")
	user := user.New(first, last, email, password)
	err := cc.DB.InsertUser(&user)
	if err != nil {
		return err
	}
	return cc.String(http.StatusOK, "ok")
}
