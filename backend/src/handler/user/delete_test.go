package user

import (
	"app/model"
	"app/response"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (suite *TestSuite) TestDelete() {

	suite.Run("normal", func() {
		suite.SetupTest()
		suite.e.POST("/user/delete", suite.h.Delete, suite.m.RequireUserSignIn)
		suite.e.POST("/user/signin", suite.h.SignIn)

		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{1})

		cookie := suite.GetSignInCookie(expected_user.ID)
		password := "12345678abc"
		post_data, _ := json.Marshal(map[string]interface{}{"password": password})

		req := httptest.NewRequest(http.MethodPost, "/user/delete", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderCookie, cookie)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)
		assert.Empty(suite.T(), res.Data)
		assert.NotEmpty(suite.T(), rec.Header().Get("Set-Cookie"))

		// session has been expired.
		req = httptest.NewRequest(http.MethodPost, "/user/delete", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderCookie, cookie)
		rec = httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)
		res = response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 401, rec.Code)

		// user has been deleted.
		var actual_user model.User
		ret := suite.seeder.DB.First(&actual_user, expected_user.ID)
		assert.ErrorIs(suite.T(), gorm.ErrRecordNotFound, ret.Error)

		suite.TearDownTest()
	})

	suite.Run("abnormal.different_password", func() {
		suite.SetupTest()
		suite.e.POST("/user/delete", suite.h.Delete, suite.m.RequireUserSignIn)
		suite.e.POST("/user/signin", suite.h.SignIn)

		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{1})

		password := "12345678abcd"
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": password})
		req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 400, rec.Code)
		assert.NotEmpty(suite.T(), res.Error)

		// user has not been deleted.
		var actual_user model.User
		ret := suite.seeder.DB.First(&actual_user, expected_user.ID)
		assert.Nil(suite.T(), ret.Error)
		assert.Equal(suite.T(), expected_user, actual_user)

		suite.TearDownTest()
	})
}
