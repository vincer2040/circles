package circles

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vincer2040/circles/internal/db"
	"github.com/vincer2040/circles/internal/render"
	"github.com/vincer2040/circles/internal/routes"
	"github.com/vincer2040/circles/internal/util"
)

func Main() {

	url, err := util.DbURL()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get dburl %s\n", err)
		os.Exit(1)
	}

	circlesDB, err := db.New(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open database %s %s\n", url, err)
		os.Exit(1)
	}

    defer circlesDB.Close()

    err = circlesDB.CreateUserTable()
    if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create user table %s\n", err)
		os.Exit(1)
    }

	e := echo.New()

	e.Renderer = render.New()

	e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &util.CirclesContext{
				Context: c,
				DB:      *circlesDB,
			}
			return next(cc)
		}
	})

	e.Static("/css", "public/styles")

	e.GET("/", routes.RootGet)
	e.GET("/signup", routes.SignupGet)
	e.POST("/signup", routes.SignupPost)

	e.Logger.Fatal(e.Start(":6969"))
}
