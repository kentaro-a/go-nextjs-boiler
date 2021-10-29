package middleware

import (
	app_context "app/context"
	"app/model"
	"app/response"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestRequireSignIn() {

	suite.e.POST("/user/signin", suite.handlers.UserHandler.SignIn)
	suite.e.POST("/user/dashboard", func(c echo.Context) error {
		cc := c.(*app_context.Context)
		return response.Success(c, 200, map[string]interface{}{"user": cc.User}, nil)
	}, suite.m.RequireUserSignIn)

	suite.Run("normal", func() {
		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{1})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": "12345678abc"})
		req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		cookie := rec.Header().Get("Set-Cookie")

		req = httptest.NewRequest(http.MethodPost, "/user/dashboard", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Cookie", cookie)
		rec = httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)

	})

	suite.Run("abnormal.invalid_session", func() {
		cookie := "mapp=MTYzMzA3MDQ0MHxCQXdBQVRJPXwReJY1tkoCaIDBycrcKTa8n3tJHMKieKhuuOjrLrpDiQ==; Path=/; Expires=Fri, 01 Oct 2021 07:40:40 GMT; Max-Age=3600"

		req := httptest.NewRequest(http.MethodPost, "/user/dashboard", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Cookie", cookie)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 401, rec.Code)
	})

	suite.Run("abnormal.missed_cookie", func() {
		req := httptest.NewRequest(http.MethodPost, "/user/dashboard", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 401, rec.Code)
	})

}
