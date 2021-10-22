package handler

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) SignUpInitial(c echo.Context) error {

	fmt.Println(c.Get("user_mail_auth"))
	return c.JSONPretty(200, map[string]string{"method": http.MethodPost}, "\t")

}

func (h Handler) SignUp(c echo.Context) error {

	fmt.Println(c.Get("user_mail_auth"))
	return c.JSONPretty(200, map[string]string{"method": http.MethodPost}, "\t")

}
