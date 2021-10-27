package handler

import (
	"app/model"
	"app/response"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestForgotPasswordVerifyToken(t *testing.T) {
	e, h, m, seeder := setup(t)
	e.POST("/forgot_password_verify_token/:token", h.ForgotPasswordVerifyToken, m.VerifyUserMailAuth)

	// 正常
	{
		token := "pre_forgot_password_ValidToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/forgot_password_verify_token/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)
	}

	teardown(t, e, seeder)
}

func TestForgotPassword(t *testing.T) {

	// 正常
	{
		e, h, m, seeder := setup(t)
		e.POST("/forgot_password/:token", h.ForgotPassword, m.VerifyUserMailAuth)

		var expected_user_mail_auth model.UserMailAuth
		seeder.DB.Find(&expected_user_mail_auth, []int64{4})

		model_user := model.NewUserModel(seeder.DB)
		prev_user, err := model_user.FindByMail(expected_user_mail_auth.Mail)
		assert.Nil(t, err)
		assert.NotEmpty(t, prev_user)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/forgot_password/%s", expected_user_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)

		model_user_mail_auth := model.NewUserMailAuthModel(seeder.DB)
		user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

		user, err := model_user.FindByMail(expected_user_mail_auth.Mail)
		assert.Nil(t, err)
		assert.NotEmpty(t, user)

		// パスワードがランダム文字列で変更されている
		assert.NotEqual(t, prev_user.Password, user.Password)

		teardown(t, e, seeder)
	}

}
