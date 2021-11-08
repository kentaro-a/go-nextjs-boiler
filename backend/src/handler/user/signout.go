package user

import (
	"app/config"
	"app/response"
	"app/session"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) SignOut(c echo.Context) error {
	// セッション削除
	session.DeleteSession(c, h.DB, config.Get().Session.Key)
	return response.Success(c, 200, nil, nil)
}
