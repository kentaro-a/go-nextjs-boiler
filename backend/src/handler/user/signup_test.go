package user

import (
	"app/model"
	"app/response"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSignUpVerifyToken(t *testing.T) {
	e, h, m, seeder := setup(t)
	e.POST("/user/signup_verify_token/:token", h.SignUpVerifyToken, m.VerifyMailAuth)

	// 正常
	{
		token := "pre_signup_ValidToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup_verify_token/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)
	}

	teardown(t, e, seeder)
}

func TestSignUp(t *testing.T) {

	// 正常
	{
		e, h, m, seeder := setup(t)
		e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

		var expected_user_mail_auth model.MailAuth
		seeder.DB.Find(&expected_user_mail_auth, []int64{1})
		name := "myname"
		password := "12345678"
		post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)

		model_user_mail_auth := model.NewMailAuthModel(seeder.DB)
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
			{
				e, h, m, seeder := setup(t)
				e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

				var expected_user_mail_auth model.MailAuth
				seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := ""
				password := "12345678"
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(t, 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				assert.Empty(t, user)

				teardown(t, e, seeder)
			}
			{
				e, h, m, seeder := setup(t)
				e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

				var expected_user_mail_auth model.MailAuth
				seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				password := "12345678"
				post_data, _ := json.Marshal(map[string]interface{}{"password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(t, 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				assert.Empty(t, user)

				teardown(t, e, seeder)
			}
			{
				e, h, m, seeder := setup(t)
				e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

				var expected_user_mail_auth model.MailAuth
				seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := strings.Repeat("1", 256)
				password := "12345678"
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(t, 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				assert.Empty(t, user)

				teardown(t, e, seeder)
			}
		}

		// password
		{
			{
				e, h, m, seeder := setup(t)
				e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

				var expected_user_mail_auth model.MailAuth
				seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				password := ""
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(t, 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				assert.Empty(t, user)

				teardown(t, e, seeder)
			}
			{
				e, h, m, seeder := setup(t)
				e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

				var expected_user_mail_auth model.MailAuth
				seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				post_data, _ := json.Marshal(map[string]interface{}{"name": name})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(t, 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(user_mail_auths))

				teardown(t, e, seeder)
			}
			{
				e, h, m, seeder := setup(t)
				e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

				var expected_user_mail_auth model.MailAuth
				seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				password := strings.Repeat("1", 65)
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(t, 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				assert.Empty(t, user)

				teardown(t, e, seeder)
			}
			{
				e, h, m, seeder := setup(t)
				e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

				var expected_user_mail_auth model.MailAuth
				seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				password := strings.Repeat("1", 7)
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(t, 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(t, err)
				assert.Equal(t, 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				assert.Empty(t, user)

				teardown(t, e, seeder)
			}
		}
	}

	// すでに登録済みのメアド(なんらかの方法で別経路で登録された場合など)
	{
		{
			e, h, m, seeder := setup(t)
			e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

			var expected_user_mail_auth model.MailAuth
			seeder.DB.Find(&expected_user_mail_auth, []int64{1})
			mail := expected_user_mail_auth.Mail
			name := "myname"
			password := strings.Repeat("1", 8)

			// ユーザーを物理的にinsertする
			user_model := model.NewUserModel(seeder.DB)
			user_model.Create(&model.User{
				Name:      "other name",
				Mail:      mail,
				Password:  "12345678",
				StatusFlg: 0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})

			post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseSuccess{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)

			teardown(t, e, seeder)
		}

		{
			e, h, m, seeder := setup(t)
			e.POST("/user/signup/:token", h.SignUp, m.VerifyMailAuth)

			var expected_user_mail_auth model.MailAuth
			seeder.DB.Find(&expected_user_mail_auth, []int64{1})
			name := "myname"
			password := strings.Repeat("1", 8)

			post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

			loop := 10
			ch := make(chan int, loop)
			for i := 0; i < loop; i++ {
				go func(post_data []byte, ch chan int) {
					req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rec := httptest.NewRecorder()
					e.ServeHTTP(rec, req)
					ch <- rec.Code

				}(post_data, ch)
			}

			code200_count := 0
			for i := 0; i < loop; i++ {
				select {
				case code := <-ch:
					if code == 200 {
						code200_count++
					}
				}
			}
			assert.Equal(t, 1, code200_count)

			teardown(t, e, seeder)
		}
	}

}
