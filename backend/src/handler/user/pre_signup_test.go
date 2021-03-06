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

func (suite *TestSuite) TestPreSignUp() {
	suite.Run("normal", func() {
		suite.SetupTest()
		suite.e.POST("/user/pre_signup", suite.h.PreSignUp)

		mail := "test@test.com"
		post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)
		assert.Empty(suite.T(), res.Data)

		m := model.NewMailAuthModel(suite.seeder.DB)
		mail_auths, err := m.FindByMailFunction(mail, "user/pre_signup")
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 1, len(mail_auths)) // Means any other records which have the same mail and function are not found.
		assert.Equal(suite.T(), mail, mail_auths[0].Mail)

		suite.TearDownTest()
	})

	suite.Run("normal.removed_old_pre_signup_records", func() {
		suite.SetupTest()
		suite.e.POST("/user/pre_signup", suite.h.PreSignUp)

		mail := "pre_signup_user1@test.com"

		// 既存の仮登録データが存在すること
		m := model.NewMailAuthModel(suite.seeder.DB)
		mail_auths_pre, err := m.FindByMailFunction(mail, "user/pre_signup")
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 1, len(mail_auths_pre)) // Means any other records which have the same mail and function are not found.

		post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)
		assert.Empty(suite.T(), res.Data)

		mail_auths, err := m.FindByMailFunction(mail, "user/pre_signup")
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 1, len(mail_auths)) // Means any other records which have the same mail and function are not found.

		suite.TearDownTest()
	})

	suite.Run("abnormal.exist_mail", func() {
		suite.SetupTest()
		suite.e.POST("/user/pre_signup", suite.h.PreSignUp)

		mail := "user1@test.com"
		post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
		req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 400, rec.Code)
		assert.Equal(suite.T(), 1, len(res.Error.Messages))

		m := model.NewMailAuthModel(suite.seeder.DB)
		mail_auths, err := m.FindByMailFunction(mail, "user/pre_signup")
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 0, len(mail_auths))

		suite.TearDownTest()
	})

	suite.Run("abnormal.invalid_mail", func() {
		suite.SetupTest()
		suite.e.POST("/user/pre_signup", suite.h.PreSignUp)

		{
			mail := "usertest.com"
			post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Equal(suite.T(), 1, len(res.Error.Messages))

			m := model.NewMailAuthModel(suite.seeder.DB)
			mail_auths, err := m.FindByMailFunction(mail, "user/pre_signup")
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), 0, len(mail_auths))
		}
		{
			mail := ""
			post_data, _ := json.Marshal(map[string]interface{}{"mail": mail})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Equal(suite.T(), 1, len(res.Error.Messages))

			m := model.NewMailAuthModel(suite.seeder.DB)
			mail_auths, err := m.FindByMailFunction(mail, "user/pre_signup")
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), 0, len(mail_auths))
		}
		{
			post_data, _ := json.Marshal(map[string]interface{}{})
			req := httptest.NewRequest(http.MethodPost, "/user/pre_signup", bytes.NewReader(post_data))
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
