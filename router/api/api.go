package api


import (
"github.com/labstack/echo"
mw "github.com/labstack/echo/middleware"
. "../../config"
)


func NewContext() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &Context{c}
			return h(ctx)
		}
	}
}

type Context struct {
	echo.Context
}

func Router() *echo.Echo {
	e := echo.New()
	e.Use()

	// Customization
	if Config.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("api")
	e.Logger.SetLevel(GetLogLvl())

	// CSRF
	e.Use(mw.CSRFWithConfig(mw.CSRFConfig{
		TokenLookup: "form:X-XSRF-TOKEN",
	}))

	// Gzip
	e.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Static("/favicon.ico", "./assets/img/favicon.ico")

	// Routers
	e.GET("/login", LoginHandler)

	// JWT
	r := e.Group("")
	r.Use(mw.JWTWithConfig(mw.JWTConfig{
		SigningKey:  []byte("secret"),
		ContextKey:  "_user",
		TokenLookup: "header:" + echo.HeaderAuthorization,
	}))

	r.GET("/", handler(ApiHandler))
	r.GET("/user", UserHandler)

	return e
}

type (
	HandlerFunc func(*Context) error
)

func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}

