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
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (suite *TestSuite) TestSignUpVerifyToken() {
	suite.e.POST("/user/signup_verify_token/:token", suite.h.SignUpVerifyToken, suite.m.VerifyMailAuth)

	// 正常
	{
		token := "pre_signup_ValidToken"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup_verify_token/%s", token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)
	}
}

func (suite *TestSuite) TestSignUp() {

	// 正常
	{
		suite.SetupTest()
		suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)

		var expected_user_mail_auth model.MailAuth
		suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
		name := "myname"
		password := "12345678"
		post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
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

		model_user := model.NewUserModel(suite.seeder.DB)
		user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
		assert.Nil(suite.T(), err)
		assert.NotEmpty(suite.T(), user)

		suite.TearDownTest()
	}

	// Error: validaion
	{
		// name
		{
			{
				suite.SetupTest()

				suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)
				var expected_user_mail_auth model.MailAuth
				suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := ""
				password := "12345678"
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				suite.e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(suite.T(), 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(suite.seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(suite.T(), err)
				assert.Equal(suite.T(), 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(suite.seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
				assert.Empty(suite.T(), user)

				suite.TearDownTest()
			}
			{
				suite.SetupTest()

				suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)
				var expected_user_mail_auth model.MailAuth
				suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				password := "12345678"
				post_data, _ := json.Marshal(map[string]interface{}{"password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				suite.e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(suite.T(), 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(suite.seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(suite.T(), err)
				assert.Equal(suite.T(), 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(suite.seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
				assert.Empty(suite.T(), user)

				suite.TearDownTest()
			}
			{
				suite.SetupTest()

				suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)
				var expected_user_mail_auth model.MailAuth
				suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := strings.Repeat("1", 256)
				password := "12345678"
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				suite.e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(suite.T(), 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(suite.seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(suite.T(), err)
				assert.Equal(suite.T(), 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(suite.seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
				assert.Empty(suite.T(), user)

				suite.TearDownTest()
			}
		}

		// password
		{
			{
				suite.SetupTest()

				suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)
				var expected_user_mail_auth model.MailAuth
				suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				password := ""
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				suite.e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(suite.T(), 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(suite.seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(suite.T(), err)
				assert.Equal(suite.T(), 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(suite.seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
				assert.Empty(suite.T(), user)

				suite.TearDownTest()
			}
			{
				suite.SetupTest()

				suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)
				var expected_user_mail_auth model.MailAuth
				suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				post_data, _ := json.Marshal(map[string]interface{}{"name": name})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				suite.e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(suite.T(), 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(suite.seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(suite.T(), err)
				assert.Equal(suite.T(), 1, len(user_mail_auths))

				suite.TearDownTest()
			}
			{
				suite.SetupTest()

				suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)
				var expected_user_mail_auth model.MailAuth
				suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				password := strings.Repeat("1", 65)
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				suite.e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(suite.T(), 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(suite.seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(suite.T(), err)
				assert.Equal(suite.T(), 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(suite.seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
				assert.Empty(suite.T(), user)

				suite.TearDownTest()
			}
			{
				suite.SetupTest()

				suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)
				var expected_user_mail_auth model.MailAuth
				suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
				name := "myname"
				password := strings.Repeat("1", 7)
				post_data, _ := json.Marshal(map[string]interface{}{"name": name, "password": password})

				req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/signup/%s", expected_user_mail_auth.Token), bytes.NewReader(post_data))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				suite.e.ServeHTTP(rec, req)

				res := response.ResponseSuccess{}
				json.NewDecoder(rec.Body).Decode(&res)
				assert.Equal(suite.T(), 400, rec.Code)

				model_user_mail_auth := model.NewMailAuthModel(suite.seeder.DB)
				user_mail_auths, err := model_user_mail_auth.FindByMailFunction(expected_user_mail_auth.Mail, expected_user_mail_auth.Function)
				assert.Nil(suite.T(), err)
				assert.Equal(suite.T(), 1, len(user_mail_auths)) // Means user_mail_auths record has been deleted.

				model_user := model.NewUserModel(suite.seeder.DB)
				user, err := model_user.FindByMailPassword(expected_user_mail_auth.Mail, model_user.GetHashedPassword(password))
				assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
				assert.Empty(suite.T(), user)

				suite.TearDownTest()
			}
		}
	}

	// すでに登録済みのメアド(なんらかの方法で別経路で登録された場合など)
	{
		{
			suite.SetupTest()

			suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)
			var expected_user_mail_auth model.MailAuth
			suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
			mail := expected_user_mail_auth.Mail
			name := "myname"
			password := strings.Repeat("1", 8)

			// ユーザーを物理的にinsertする
			user_model := model.NewUserModel(suite.seeder.DB)
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
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseSuccess{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)

			suite.TearDownTest()
		}

		{
			suite.SetupTest()
			suite.e.POST("/user/signup/:token", suite.h.SignUp, suite.m.VerifyMailAuth)

			var expected_user_mail_auth model.MailAuth
			suite.seeder.DB.Find(&expected_user_mail_auth, []int64{1})
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
					suite.e.ServeHTTP(rec, req)
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
			assert.Equal(suite.T(), 1, code200_count)

			suite.TearDownTest()
		}
	}

}
