package middleware

import (
	"app/model"
	"app/response"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestVerifyMailAuth() {
	suite.e.POST("/test/:token", func(c echo.Context) error {
		return response.Success(c, 200, nil, nil)
	}, suite.m.VerifyMailAuth)

	suite.Run("normal", func() {
		var expected_mail_auth model.MailAuth
		suite.seeder.DB.Find(&expected_mail_auth, []int64{1})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/test/%s", expected_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)
	})

	suite.Run("abnormal.expired_token", func() {
		var expected_mail_auth model.MailAuth
		suite.seeder.DB.Find(&expected_mail_auth, []int64{2})
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/test/%s", expected_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 401, rec.Code)

	})

	suite.Run("abnormal.inactive_token", func() {
		var expected_mail_auth model.MailAuth
		suite.seeder.DB.Find(&expected_mail_auth, []int64{3})
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/test/%s", expected_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 401, rec.Code)
	})

}
