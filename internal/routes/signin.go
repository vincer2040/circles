package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/util"
)

func SigninGet(c echo.Context) error {
	cc := c.(*util.CirclesContext)
	session, _ := cc.Store.Get(c.Request(), "auth")
	if !session.IsNew {
		authed, ok := session.Values["authenticated"].(bool)
		if ok && authed {
			c.Redirect(http.StatusSeeOther, "/home")
		}
	}

	pageData := util.PageData{
		Route: "/signin",
	}

	return cc.Render(http.StatusOK, "signin.html", pageData)
}

func SigninPost(c echo.Context) error {
	cc := c.(*util.CirclesContext)
	email := cc.FormValue("email")
	password := cc.FormValue("password")
	user, err := cc.DB.GetUser(email)
	if err != nil {
		return err
	}
	authed := user.Authenticate(password)
	if !authed {
		return cc.HTML(http.StatusUnauthorized, "failed to Authenticate")
	}
	session, _ := cc.Store.Get(c.Request(), "auth")
	session.Values["authenticated"] = true
	session.Values["first"] = user.First
	session.Values["email"] = user.Email
	session.Save(cc.Request(), cc.Response())
    cc.Response().Header().Add("Hx-Redirect", "/home")
	return cc.String(http.StatusOK, "")
}
