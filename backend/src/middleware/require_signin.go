package middleware

import (
	app_context "app/context"
	"app/response"
	app_session "app/session"

	echo "github.com/labstack/echo/v4"
)

func (m Middleware) RequireSignIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		is_signedin, user := app_session.IsSignedIn(c)
		if !is_signedin {
			return response.Error(c, 401, nil, nil)
		}
		cc := c.(*app_context.Context)
		cc.User = &user
		return next(cc)
	}
}
