package router

import (
	"app/handler"
	app_middleware "app/middleware"
	"app/model"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func New() (*echo.Echo, error) {

	db, err := model.NewDB()
	if err != nil {
		return nil, err
	}

	h := handler.Handler{DB: db}
	m := app_middleware.Middleware{DB: db}
	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(m.Session())
	e.Use(app_middleware.Cors())
	e.Use(app_middleware.AccessLog)

	// handlers
	e.POST("/signin", h.SignIn)
	e.POST("/pre_signup", h.PreSignUp)
	e.POST("/signup_verify_token/:token", h.SignUpVerifyToken, m.VerifyUserMailAuth)
	e.POST("/signup/:token", h.SignUp, m.VerifyUserMailAuth)

	e.GET("/dashboard", h.Dashboard, m.UserAuthentication)

	return e, nil
}
