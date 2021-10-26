package handler

import (
	"app/model"
	"app/response"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSignUpVerifyToken(t *testing.T) {
	e, h, m, seeder := setup(t)
	e.POST("/signup_verify_token/:token", h.SignUpVerifyToken, m.VerifyUserMailAuth)

	// 正常
	{
		token := "pre_signup_ValidToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup_verify_token/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)
	}
	// Error: Expired token
	{
		token := "pre_signup_ExpiredToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup_verify_token/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 401, rec.Code)
	}
	// Error: InAcrive token
	{
		token := "pre_signup_InActiveToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup_verify_token/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 401, rec.Code)
	}
	// Error: Without token
	{
		token := ""
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup_verify_token/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 404, rec.Code)
	}
	teardown(t, e, seeder)
}

func TestSignUp(t *testing.T) {

	// Valid Token
	{
		// 正常
		{
			e, h, m, seeder := setup(t)
			e.POST("/signup/:token", h.SignUp, m.VerifyUserMailAuth)

			var expected_user_mail_auth model.UserMailAuth
			seeder.DB.Find(&expected_user_mail_auth, []int64{1})
			name := "myname"
			password := "12345678"
			post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
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

			model_user := model.NewUserModel(seeder.DB)
			user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
			assert.Nil(t, err)
			assert.NotEmpty(t, user)

			teardown(t, e, seeder)
		}

		// Error: validaion
		{
			// name
			{
				e, h, m, seeder := setup(t)
				e.POST("/signup/:token", h.SignUp, m.VerifyUserMailAuth)

				var expected_user_mail_auth model.UserMailAuth
				seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				password := "12345678"
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
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

				model_user := model.NewUserModel(seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.Nil(t, err)
				assert.NotEmpty(t, user)

				teardown(t, e, seeder)
			}

		}

	}

	// Error: Expired token
	{

		e, h, m, seeder := setup(t)
		e.POST("/signup/:token", h.SignUp, m.VerifyUserMailAuth)

		token := "pre_signup_ExpiredToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 401, rec.Code)

		teardown(t, e, seeder)
	}
	// Error: InAcrive token
	{
		e, h, m, seeder := setup(t)
		e.POST("/signup/:token", h.SignUp, m.VerifyUserMailAuth)

		token := "pre_signup_InActiveToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 401, rec.Code)
		teardown(t, e, seeder)
	}
	// Error: Without token
	{
		e, h, m, seeder := setup(t)
		e.POST("/signup/:token", h.SignUp, m.VerifyUserMailAuth)

		token := ""
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/signup/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 404, rec.Code)

		teardown(t, e, seeder)
	}

}
