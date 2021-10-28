package user

import (
	"app/response"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	return response.Success(c, 200, nil, nil)
}
