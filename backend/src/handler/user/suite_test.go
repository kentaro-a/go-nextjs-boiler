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

	"github.com/stretchr/testify/suite"

	middleware "github.com/labstack/echo/v4/middleware"

	echo "github.com/labstack/echo/v4"
)

type TestSuite struct {
	suite.Suite
	e      *echo.Echo
	h      Handler
	m      app_middleware.Middleware
	seeder *tests.Seeder
}

func (suite *TestSuite) SetupTest() {
	suite.e = echo.New()
	suite.e.Use(app_middleware.Context)
	suite.e.Use(middleware.Logger())
	suite.e.Use(middleware.Recover())
	suite.e.Use(app_middleware.Cors())
	suite.e.Use(app_middleware.AccessLog)

	suite.seeder = tests.NewSeeder()
	unSeedAll(suite.seeder)
	suite.seeder.Seed(tests.MailAuthsFixture(), tests.UsersFixture())
	suite.m = app_middleware.Middleware{DB: suite.seeder.DB}
	suite.e.Use(suite.m.Session())
	suite.h = Handler{DB: suite.seeder.DB}
}

func (suite *TestSuite) TearDownTest() {
	app_session.DeleteStore()
	// unSeedAll(suite.seeder)
	suite.seeder.Close()
}

func unSeedAll(seeder *tests.Seeder) {
	truncate_tables := []string{
		"users",
		"mail_auths",
		"sessions",
	}
	seeder.UnSeed(truncate_tables...)
}

func (suite *TestSuite) GetSignInCookie(user_id uint) string {
	suite.e.POST("/user/signin", suite.h.SignIn)
	var user model.User
	suite.seeder.DB.First(&user, user_id)
	post_data, _ := json.Marshal(map[string]interface{}{"mail": user.Mail, "password": "12345678abc"})
	req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	suite.e.ServeHTTP(rec, req)
	cookie := rec.Header().Get("Set-Cookie")
	return cookie
}

func (suite *TestSuite) GetInvalidSignInCookie() string {
	return "mapp=MTYzNTI5ODI3M3xCQXdBQVRFPXwQDLpnxxtZT-ETDMY3pfj5tH3OQzjioupuQ8G0o45e4w==; Path=/; Expires=Wed, 27 Oct 2021 02:31:13 GMT; Max-Age=3600"
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
