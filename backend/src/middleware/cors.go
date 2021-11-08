package middleware

import (
	"app/config"
	"net/http"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func Cors() echo.MiddlewareFunc {
	c := config.Get()
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: c.Web.Cors.AllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},
	})

}
