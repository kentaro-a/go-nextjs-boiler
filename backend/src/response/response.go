package response

import (
	app_error "app/error"
	app_log "app/log"

	echo "github.com/labstack/echo/v4"
)

type ResponseSuccess struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

func Success(c echo.Context, status_code int, data interface{}, log_fields *app_log.Fields) error {
	res := ResponseSuccess{
		status_code,
		data,
	}

	if log_fields != nil {
		app_log.Error(c, *log_fields)
	}
	return c.JSONPretty(status_code, res, "\t")
}

type ResponseError struct {
	Error app_error.Error `json:"error"`
}

func Error(c echo.Context, status_code int, messages []string, log_fields *app_log.Fields) error {
	res := ResponseError{
		app_error.New(status_code, messages),
	}
	if log_fields != nil {
		app_log.Error(c, *log_fields)
	}
	return c.JSONPretty(status_code, res, "\t")
}

func SystemError(c echo.Context, log_fields *app_log.Fields) error {
	res := ResponseError{
		app_error.New(500, []string{"予期せぬエラーが発生しました。時間をおいてから再度お試しください。"}),
	}
	if log_fields != nil {
		app_log.Error(c, *log_fields)
	}
	return c.JSONPretty(500, res, "\t")
}
