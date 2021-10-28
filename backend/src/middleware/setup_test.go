package middleware

import (
	user_handler "app/handler/user"
	app_session "app/session"
	"app/tests"
	"testing"

	middleware "github.com/labstack/echo/v4/middleware"

	echo "github.com/labstack/echo/v4"
)

type Handlers struct {
	UserHandler user_handler.Handler
}

func setup(t *testing.T) (*echo.Echo, Handlers, Middleware, *tests.Seeder) {
	e := echo.New()
	e.Use(Context)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(Cors())
	e.Use(AccessLog)

	seeder := tests.NewSeeder()
	seeder.Seed(tests.MailAuthsFixture(), tests.UsersFixture())
	m := Middleware{DB: seeder.DB}
	e.Use(m.Session())
	h := user_handler.Handler{DB: seeder.DB}

	handlers := Handlers{
		UserHandler: h,
	}
	return e, handlers, m, seeder

}

func teardown(t *testing.T, e *echo.Echo, seeder *tests.Seeder) {
	app_session.DeleteStore()
	truncate_tables := []string{
		"users",
		"mail_auths",
		"sessions",
	}
	seeder.UnSeed(truncate_tables...)
	seeder.Close()
}
