package router

import (
	user_handler "app/handler/user"
	app_middleware "app/middleware"
	"app/model"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func New() (*echo.Echo, error) {

	e := echo.New()

	db, err := model.NewDB()
	if err != nil {
		return nil, err
	}

	m := app_middleware.Middleware{DB: db}

	// middlewares
	e.Use(app_middleware.Context)
	e.Use(middleware.Recover())
	e.Use(m.Session())
	e.Use(app_middleware.Cors())
	e.Use(app_middleware.AccessLog)

	{
		g := e.Group("/user")
		h := user_handler.Handler{DB: db}

		g.POST("/pre_signup", h.PreSignUp)
		g.POST("/signup_verify_token/:token", h.SignUpVerifyToken, m.VerifyMailAuth)
		g.POST("/signup/:token", h.SignUp, m.VerifyMailAuth)
		g.POST("/signin", h.SignIn)
		g.POST("/signout", h.SignOut, m.RequireUserSignIn)
		g.POST("/pre_forgot_password", h.PreForgotPassword)
		g.POST("/forgot_password_verify_token/:token", h.ForgotPasswordVerifyToken, m.VerifyMailAuth)
		g.POST("/forgot_password/:token", h.ForgotPassword, m.VerifyMailAuth)
		g.POST("/delete", h.Delete, m.RequireUserSignIn)
		g.POST("/dashboard", h.Dashboard, m.RequireUserSignIn)
	}

	return e, nil
}
