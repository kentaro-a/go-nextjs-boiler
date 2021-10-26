package context

import (
	"app/model"

	echo "github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	User   *model.User `json:"user"`
	Status *Status     `json:"status"`
}

type Status struct {
	BidStatusFlg  int `json:"bid_status_flg"`
	LineStatusFlg int `json:"line_status_flg"`
}

func NewContext(c echo.Context) *Context {
	return &Context{
		c,
		nil,
		nil,
	}
}
