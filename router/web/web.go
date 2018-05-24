package web

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
	e.Use(NewContext())

	// Customization
	if Config.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("web")
	e.Logger.SetLevel(GetLogLvl())


	// CSRF
	e.Use(mw.CSRFWithConfig(mw.CSRFConfig{
		ContextKey:  "_csrf",
		TokenLookup: "form:_csrf",
	}))

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	// Routers
	e.GET("/", handler(HomeHandler))

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