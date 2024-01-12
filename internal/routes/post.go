package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

    image := cc.FormValue("image")
    description := cc.FormValue("description")

    now := time.Now().UTC()

    timeString := now.Format("2006-01-02T15:04:05Z")

    // TODO: save the file to disc and path to file in db with post
    fileName := fmt.Sprintf("images/%s/%s/%s", circle, author, timeString)
    fmt.Println("filename", fileName, "image", image)

    err := cc.DB.InsertPost(circle, author, description, timeString)
    if err != nil {
        return err
    }

    return cc.String(http.StatusOK, "")
}

func PostDelete(c echo.Context) error {
    cc := c.(*util.CirclesContext)

    session, err := cc.Store.Get(c.Request(), "auth")
    if err != nil {
        return cc.Redirect(http.StatusSeeOther, "/signin")
    }

    if session.IsNew {
        return cc.Redirect(http.StatusSeeOther, "/signin")
    }

    _, ok := session.Values["email"].(string)

    if !ok {
        return cc.Redirect(http.StatusSeeOther, "/signin")
    }

    postIDString := cc.Param("id")
    postId, err := strconv.ParseInt(postIDString, 10, 64)
    if err != nil {
        return err
    }

    err = cc.DB.DeletePost(postId)
    if err != nil {
        return err
    }
    return cc.String(http.StatusOK, "")
}
