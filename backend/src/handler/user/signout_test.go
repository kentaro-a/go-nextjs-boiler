package user

import (
	"app/model"
	"fmt"
	"net/http"
	"net/http/httptest"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestSignOut() {
	suite.Run("normal", func() {
		suite.e.POST("/user/signout", suite.h.SignOut, suite.m.RequireUserSignIn)
		suite.e.POST("/user/test", func(c echo.Context) error {
			return c.JSON(200, nil)
		}, suite.m.RequireUserSignIn)

		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{1})
		cookie := suite.GetSignInCookie(expected_user.ID)

		// succeeded to access with valid cookie
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/test"), nil)
		req.Header.Set(echo.HeaderCookie, cookie)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)
		assert.Equal(suite.T(), 200, rec.Code)

		// signout
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signout"), nil)
		req.Header.Set(echo.HeaderCookie, cookie)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)
		assert.Equal(suite.T(), 200, rec.Code)

		// failed to access with revoked cookie
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signout"), nil)
		req.Header.Set(echo.HeaderCookie, cookie)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)
		assert.Equal(suite.T(), 401, rec.Code)
	})

	suite.Run("abnormal.invalid_cookie", func() {
		suite.e.POST("/user/signout", suite.h.SignOut, suite.m.RequireUserSignIn)

		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{1})
		cookie := suite.GetInvalidSignInCookie()

		// signout
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signout"), nil)
		req.Header.Set(echo.HeaderCookie, cookie)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)
		assert.Equal(suite.T(), 401, rec.Code)
	})

}
