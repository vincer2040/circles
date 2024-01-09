package circles

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
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

	key, err := util.GenerateRandomKey(32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate random key %s\n", err)
		os.Exit(1)
	}

	circlesDB, err := db.New(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open database %s %s\n", url, err)
		os.Exit(1)
	}

	defer circlesDB.Close()

	/*
		    err = circlesDB.DropUserTable()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to drop user table %s\n", err)
				os.Exit(1)
			}
	*/

	err = circlesDB.CreateUserTable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create user table %s\n", err)
		os.Exit(1)
	}

	err = circlesDB.CreateCirclesTable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create circles table %s\n", err)
		os.Exit(1)
	}

	err = circlesDB.CreateCircleUsersTable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create circles users table %s\n", err)
		os.Exit(1)
	}

	err = circlesDB.CreatePostsTable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create posts table %s\n", err)
		os.Exit(1)
	}

	e := echo.New()

	e.Renderer = render.New()

	store := sessions.NewCookieStore([]byte(key))

	e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &util.CirclesContext{
				Context: c,
				DB:      *circlesDB,
				Store:   store,
			}
			return next(cc)
		}
	})

	e.Static("/css", "public/styles")

	e.GET("/", routes.RootGet)
	e.GET("/signup", routes.SignupGet)
	e.POST("/signup", routes.SignupPost)
	e.GET("/signin", routes.SigninGet)
	e.POST("/signin", routes.SigninPost)
	e.GET("/me", routes.MeGet)
	e.GET("/home", routes.HomeGet)
	e.GET("/create-circle", routes.CreateCircleGet)
	e.POST("/create-circle", routes.CreateCirclePost)
	e.GET("/circle/:circle", routes.CircleGet)
    e.POST("/post/circle/:circle", routes.PostPost)

	e.Logger.Fatal(e.Start(":6969"))
}
