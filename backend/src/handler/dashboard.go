package handler

import (
	app_context "app/context"
	"app/response"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) Dashboard(c echo.Context) error {
	cc := app_context.CastContext(c)
	return response.Success(c, 200, map[string]interface{}{"user": cc.User}, nil)
}
