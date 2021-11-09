package user

import (
	app_context "app/context"
	"app/response"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) VerifySignIn(c echo.Context) error {
	// ミドルウェアでログインセッションは確認済み
	cc := app_context.CastContext(c)
	return response.Success(c, 200, map[string]interface{}{
		"user": cc.User,
	}, nil)
}
