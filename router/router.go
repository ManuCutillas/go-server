package router

import (
"context"
"net/url"
"os"
"os/signal"
"time"

"github.com/labstack/echo"
mw "github.com/labstack/echo/middleware"
"github.com/labstack/gommon/log"

. "../config"
"./api"
"./socket"
"./web"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func InitRoutes() map[string]*Host {
	// Hosts
	hosts := make(map[string]*Host)

	hosts[Config.Server.DomainWeb] = &Host{web.Router()}
	hosts[Config.Server.DomainApi] = &Host{api.Router()}
	hosts[Config.Server.DomainSocket] = &Host{socket.Router()}

	return hosts
}

// Subdomain deployment
func RunSubdomains(confFilePath string) {
	// Configure initialization
	if err := InitConfig(confFilePath); err != nil {
		log.Panic(err)
	}

	// Global log level
	log.SetLevel(GetLogLvl())

	// Server
	e := echo.New()
	e.Pre(mw.RemoveTrailingSlash())

	// log level
	e.Logger.SetLevel(GetLogLvl())

	// Secure, XSS/CSS HSTS
	e.Use(mw.SecureWithConfig(mw.DefaultSecureConfig))
	mw.MethodOverride()

	// CORS
	e.Use(mw.CORSWithConfig(mw.CORSConfig{
		AllowOrigins: []string{"http://" + Config.Server.DomainWeb, "http://" + Config.Server.DomainApi},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding, echo.HeaderAuthorization},
	}))

	hosts := InitRoutes()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()

		u, _err := url.Parse(c.Scheme() + "://" + req.Host)
		if _err != nil {
			e.Logger.Errorf("Request URL parse error:%v", _err)
		}

		host := hosts[u.Hostname()]
		if host == nil {
			e.Logger.Info("Host not found")
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	if !Config.Server.Graceful {
		e.Logger.Fatal(e.Start(Config.Server.Addr))
	} else {
		// Graceful Shutdown
		// Start server
		go func() {
			if err := e.Start(Config.Server.Addr); err != nil {
				e.Logger.Errorf("Shutting down the server with error:%v", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
