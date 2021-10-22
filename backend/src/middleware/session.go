package middleware

import (
	app_session "app/session"

	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
)

func (m Middleware) Session() echo.MiddlewareFunc {
	store := app_session.GetStore(m.DB)
	return session.Middleware(store)
}
