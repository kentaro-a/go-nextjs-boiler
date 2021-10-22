package middleware

import (
	app_handler "app/handler"
	app_session "app/session"
	"app/tests"
	"testing"

	middleware "github.com/labstack/echo/v4/middleware"

	echo "github.com/labstack/echo/v4"
)

func setup(t *testing.T) (*echo.Echo, app_handler.Handler, Middleware, *tests.Seeder) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(Cors())
	e.Use(AccessLog)

	seeder := tests.NewSeeder()
	seeder.Seed(tests.UserMailAuthsFixture(), tests.UsersFixture())
	m := Middleware{DB: seeder.DB}
	e.Use(m.Session())
	h := app_handler.Handler{DB: seeder.DB}
	return e, h, m, seeder

}

func teardown(t *testing.T, e *echo.Echo, seeder *tests.Seeder) {
	app_session.DeleteStore()
	truncate_tables := []string{
		"users",
		"user_mail_auths",
		"sessions",
	}
	seeder.UnSeed(truncate_tables...)
	seeder.Close()
}
