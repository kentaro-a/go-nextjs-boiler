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
)

func (suite *TestSuite) TestPreForgotPassword() {

	suite.Run("normal", func() {
		suite.SetupTest()
		suite.e.POST("/user/pre_forgot_password", suite.h.PreForgotPassword)

		mail := "user1@test.com"
		post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_forgot_password", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)

		m := model.NewMailAuthModel(suite.seeder.DB)
		user_mail_auths, err := m.FindByMailFunction(mail, "user/pre_forgot_password")
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 1, len(user_mail_auths))

		suite.TearDownTest()
	})

	suite.Run("abnormal.unexist_mail", func() {
		suite.SetupTest()
		suite.e.POST("/user/pre_forgot_password", suite.h.PreForgotPassword)

		mail := "test@test.com"
		post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_forgot_password", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 400, rec.Code)
		assert.Empty(suite.T(), res.Data)

		m := model.NewMailAuthModel(suite.seeder.DB)
		user_mail_auths, err := m.FindByMailFunction(mail, "user/pre_forgot_password")
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 0, len(user_mail_auths))

		suite.TearDownTest()
	})

	// Error: 不正なメールアドレス
	suite.Run("abnormal.invalid_mail", func() {
		suite.SetupTest()
		suite.e.POST("/user/pre_forgot_password", suite.h.PreForgotPassword)

		{
			mail := "usertest.com"
			post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_forgot_password", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Equal(suite.T(), 1, len(res.Error.Messages))

			m := model.NewMailAuthModel(suite.seeder.DB)
			user_mail_auths, err := m.FindByMailFunction(mail, "user/pre_forgot_password")
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), 0, len(user_mail_auths))
		}
		{
			mail := ""
			post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_forgot_password", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Equal(suite.T(), 1, len(res.Error.Messages))

			m := model.NewMailAuthModel(suite.seeder.DB)
			user_mail_auths, err := m.FindByMailFunction(mail, "user/pre_forgot_password")
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), 0, len(user_mail_auths))
		}
		{
			post_data, _ := json.Marshal(map[string]interface{}{})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_forgot_password", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Equal(suite.T(), 1, len(res.Error.Messages))
		}

		suite.TearDownTest()
	})
}
