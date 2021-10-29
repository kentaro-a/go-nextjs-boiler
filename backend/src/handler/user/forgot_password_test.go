package user

import (
	"app/model"
	"app/response"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestForgotPasswordVerifyToken() {
	suite.Run("normal", func() {
		suite.e.POST("/user/forgot_password_verify_token/:token",
			suite.h.ForgotPasswordVerifyToken, suite.m.VerifyMailAuth)

		token := "pre_forgot_password_ValidToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/forgot_password_verify_token/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)
	})

}

func (suite *TestSuite) TestForgotPassword() {

	suite.Run("normal", func() {
		suite.SetupTest()

		suite.e.POST("/user/forgot_password/:token", suite.h.ForgotPassword, suite.m.VerifyMailAuth)

		var expected_user_mail_auth model.MailAuth
		suite.seeder.DB.Find(&expected_user_mail_auth, []int64{4})

		model_user := model.NewUserModel(suite.seeder.DB)
		prev_user, err := model_user.FindByMail(expected_user_mail_auth.Mail)
		assert.Nil(suite.T(), err)
		assert.NotEmpty(suite.T(), prev_user)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/forgot_password/%s", expected_user_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)

		model_user_mail_auth := model.NewMailAuthModel(suite.seeder.DB)
		user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 0, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

		user, err := model_user.FindByMail(expected_user_mail_auth.Mail)
		assert.Nil(suite.T(), err)
		assert.NotEmpty(suite.T(), user)

		// パスワードがランダム文字列で変更されている
		assert.NotEqual(suite.T(), prev_user.Password, user.Password)

		suite.TearDownTest()
	})

	suite.Run("abnormal.token_expired_by_other_token", func() {
		suite.SetupTest()

		suite.e.POST("/user/forgot_password/:token", suite.h.ForgotPassword, suite.m.VerifyMailAuth)
		suite.e.POST("/user/pre_forgot_password", suite.h.PreForgotPassword)

		var expected_user_mail_auth model.MailAuth
		suite.seeder.DB.Find(&expected_user_mail_auth, []int64{4})

		model_user := model.NewUserModel(suite.seeder.DB)
		prev_user, err := model_user.FindByMail(expected_user_mail_auth.Mail)
		assert.Nil(suite.T(), err)
		assert.NotEmpty(suite.T(), prev_user)

		// publish other token.
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user_mail_auth.Mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_forgot_password", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		// expired prev token
		req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/forgot_password/%s", expected_user_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 401, rec.Code)

		suite.TearDownTest()
	})

}
