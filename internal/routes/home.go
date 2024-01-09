package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/util"
)

func HomeGet(c echo.Context) error {
	cc := c.(*util.CirclesContext)

	session, _ := cc.Store.Get(cc.Request(), "auth")
	if session.IsNew {
		return cc.Redirect(http.StatusSeeOther, "/signin")
	}

	first, ok := session.Values["first"].(string)
    email, ok := session.Values["email"].(string)

	if !ok {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

    circles, err := cc.DB.GetUsersCircles(email)
    if err != nil {
        return err
    }

    if len(circles) == 0 {
        fmt.Println("no circles")
    }

    for _, circle := range(circles) {
        fmt.Println("here:", circle)
    }

	return cc.Render(http.StatusOK, "home.html", map[string]interface{}{
		"Route": "/home",
		"First": first,
        "CirclesLength": len(circles),
        "Circles": circles,
	})
}
