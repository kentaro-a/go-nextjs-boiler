package middleware

import (
	app_log "app/log"
	"app/model"
	"app/response"
	"runtime"

	echo "github.com/labstack/echo/v4"
)

func (m Middleware) VerifyMailAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user_mail_auth := model.UserMailAuth{}
		user_mail_auth.Token = c.Param("token")
		user_mail_auth_model := model.NewUserMailAuthModel(m.DB)
		is_expired, err := user_mail_auth_model.IsExpiredToken(&user_mail_auth)
		if err != nil {
			return response.SystemError(c, &app_log.Fields{
				ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
				Messages:   []string{err.Error()},
				Error:      err,
			})
		}

		if is_expired == true {
			return response.Error(c, 401, nil, nil)
		}
		c.Set("user_mail_auth", user_mail_auth)
		return next(c)
	}
}
