package middleware

import (
	app_context "app/context"

	echo "github.com/labstack/echo/v4"
)

func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Wrap with custom context
		c = app_context.NewContext(c)
		return next(c)
	}
}
