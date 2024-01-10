package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/util"
)

func CircleGet(c echo.Context) error {
	cc := c.(*util.CirclesContext)
	session, _ := cc.Store.Get(cc.Request(), "auth")
	if session.IsNew {
		return cc.Redirect(http.StatusSeeOther, "/signin")
	}

	email, ok := session.Values["email"].(string)

	circle := cc.Param("circle")
	if !ok {
		return cc.Redirect(http.StatusSeeOther, "/signin")
	}

	isInCircle, err := cc.DB.UserIsInCircle(circle, email)
	if err != nil {
		return err
	}

	if !isInCircle {
		return cc.Redirect(http.StatusSeeOther, "/signin")
	}

	posts, err := cc.DB.GetPostsForCircle(circle)
	if err != nil {
		return err
	}

	return cc.Render(http.StatusOK, "circle.html", map[string]interface{}{
		"Route":    "/circle",
		"Circle":   circle,
		"PostsLen": len(posts),
		"Posts":    posts,
	})
}
