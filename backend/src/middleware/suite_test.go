package middleware

import (
	user_handler "app/handler/user"
	app_session "app/session"
	"app/tests"
	"testing"

	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/suite"

	echo "github.com/labstack/echo/v4"
)

type TestSuite struct {
	suite.Suite
	e        *echo.Echo
	handlers Handlers
	m        Middleware
	seeder   *tests.Seeder
}

func (suite *TestSuite) SetupTest() {
	suite.e = echo.New()
	suite.e.Use(Context)
	suite.e.Use(middleware.Logger())
	suite.e.Use(middleware.Recover())
	suite.e.Use(Cors())
	suite.e.Use(AccessLog)

	suite.seeder = tests.NewSeeder()
	unSeedAll(suite.seeder)
	suite.seeder.Seed(tests.MailAuthsFixture(), tests.UsersFixture())
	suite.m = Middleware{DB: suite.seeder.DB}
	suite.e.Use(suite.m.Session())

	h := user_handler.Handler{DB: suite.seeder.DB}
	suite.handlers = Handlers{
		UserHandler: h,
	}
}

func (suite *TestSuite) TearDownTest() {
	app_session.DeleteStore()
	unSeedAll(suite.seeder)
	suite.seeder.Close()
}

type Handlers struct {
	UserHandler user_handler.Handler
}

func unSeedAll(seeder *tests.Seeder) {
	truncate_tables := []string{
		"users",
		"mail_auths",
		"sessions",
	}
	seeder.UnSeed(truncate_tables...)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
