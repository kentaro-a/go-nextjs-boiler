package middleware

import (
	app_log "app/log"
	"runtime"

	echo "github.com/labstack/echo/v4"
)

func AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Put accesslog
		app_log.Info(c, app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
		})
		return next(c)
	}
}
