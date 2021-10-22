package handler

import (
	"app/model"
	"app/response"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) Dashboard(c echo.Context) error {
	user, _ := c.Get("user").(model.User)
	return response.Success(c, 200, map[string]interface{}{"user": user}, nil)
}
