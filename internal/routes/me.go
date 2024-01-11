package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vincer2040/circles/internal/util"
)

func MeGet(c echo.Context) error {
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

	posts, err := cc.DB.GetPostsForUser(email)

	if err != nil {
		return err
	}

	circlesCreated, err := cc.DB.GetCreatorCircles(email)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "me.html", map[string]interface{}{
		"Route":             "/me",
		"First":             first,
		"Email":             email,
		"PostsLen":          len(posts),
		"Posts":             posts,
		"CirclesCreatedLen": len(circlesCreated),
		"CirclesCreated":    circlesCreated,
	})
}
