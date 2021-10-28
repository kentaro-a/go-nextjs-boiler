package user

import (
	app_middleware "app/middleware"
	"app/model"
	app_session "app/session"
	"app/tests"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware "github.com/labstack/echo/v4/middleware"

	echo "github.com/labstack/echo/v4"
)

func setup(t *testing.T) (*echo.Echo, Handler, app_middleware.Middleware, *tests.Seeder) {
	e := echo.New()
	e.Use(app_middleware.Context)
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

func getSignInCookie(e *echo.Echo, h Handler, m app_middleware.Middleware, seeder *tests.Seeder, user_id int64) string {
	e.POST("/user/signin", h.SignIn)
	var user model.User
	seeder.DB.Find(&user, []int64{user_id})
	post_data, _ := json.Marshal(map[string]interface{}{"mail": user.Mail, "password": "12345678abc"})
	req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	cookie := rec.Header().Get("Set-Cookie")
	return cookie
}
