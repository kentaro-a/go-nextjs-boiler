package handler

import (
	app_middleware "app/middleware"
	app_session "app/session"
	"app/tests"
	"testing"

	middleware "github.com/labstack/echo/v4/middleware"

	echo "github.com/labstack/echo/v4"
)

func setup(t *testing.T) (*echo.Echo, Handler, app_middleware.Middleware, *tests.Seeder) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(app_middleware.Cors())
	e.Use(app_middleware.AccessLog)

	seeder := tests.NewSeeder()
	unSeedAll(seeder)
	seeder.Seed(tests.UserMailAuthsFixture(), tests.UsersFixture())
	m := app_middleware.Middleware{DB: seeder.DB}
	e.Use(m.Session())
	h := Handler{DB: seeder.DB}
	return e, h, m, seeder

}

func teardown(t *testing.T, e *echo.Echo, seeder *tests.Seeder) {
	app_session.DeleteStore()
	unSeedAll(seeder)
	seeder.Close()
}

func unSeedAll(seeder *tests.Seeder) {
	truncate_tables := []string{
		"users",
		"user_mail_auths",
		"sessions",
	}
	seeder.UnSeed(truncate_tables...)
}
