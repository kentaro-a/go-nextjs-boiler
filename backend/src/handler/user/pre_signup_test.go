package user

import (
	"app/model"
	"app/response"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPreSignUp(t *testing.T) {
	e, h, _, seeder := setup(t)

	e.POST("/user/pre_signup", h.PreSignUp)

	// 正常
	{
		mail := "test@test.com"
		post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)
		assert.Empty(t, res.Data)

		m := model.NewUserMailAuthModel(seeder.DB)
		user_mail_auths, err := m.FindByMailFunction(mail, "pre_signup")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(user_mail_auths)) // Means any other records which have the same mail and function are not found.
		assert.Equal(t, mail, user_mail_auths[0].Mail)
	}

	// 正常: すでに登録済みの仮登録データが削除され、最新のもののみ保持される
	{
		mail := "pre_signup_user1@test.com"

		// 既存の仮登録データが存在すること
		m := model.NewUserMailAuthModel(seeder.DB)
		user_mail_auths_pre, err := m.FindByMailFunction(mail, "pre_signup")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(user_mail_auths_pre)) // Means any other records which have the same mail and function are not found.

		post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)
		assert.Empty(t, res.Data)

		user_mail_auths, err := m.FindByMailFunction(mail, "pre_signup")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(user_mail_auths)) // Means any other records which have the same mail and function are not found.
	}

	// Error: 登録済みメールアドレス
	{
		mail := "user1@test.com"
		post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 400, rec.Code)
		assert.Equal(t, 1, len(res.Error.Messages))

		m := model.NewUserMailAuthModel(seeder.DB)
		user_mail_auths, err := m.FindByMailFunction(mail, "pre_signup")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(user_mail_auths))
	}

	// Error: 不正なメールアドレス
	{
		{
			mail := "usertest.com"
			post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Equal(t, 1, len(res.Error.Messages))

			m := model.NewUserMailAuthModel(seeder.DB)
			user_mail_auths, err := m.FindByMailFunction(mail, "pre_signup")
			assert.Nil(t, err)
			assert.Equal(t, 0, len(user_mail_auths))
		}
		{
			mail := ""
			post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Equal(t, 1, len(res.Error.Messages))

			m := model.NewUserMailAuthModel(seeder.DB)
			user_mail_auths, err := m.FindByMailFunction(mail, "pre_signup")
			assert.Nil(t, err)
			assert.Equal(t, 0, len(user_mail_auths))
		}
		{
			post_data, _ := json.Marshal(map[string]interface{}{})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Equal(t, 1, len(res.Error.Messages))
		}
	}
	teardown(t, e, seeder)

}
