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

func (suite *TestSuite) TestSignIn() {

	suite.Run("normal", func() {
		suite.SetupTest()
		suite.e.POST("/user/signin", suite.h.SignIn)

		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{1})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": "12345678abc"})
		req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 200, rec.Code)
		assert.Empty(suite.T(), res.Data)
		assert.NotEmpty(suite.T(), rec.Header().Get("Set-Cookie"))

		suite.TearDownTest()
	})

	suite.Run("abnormal.inactive_user", func() {
		suite.SetupTest()
		suite.e.POST("/user/signin", suite.h.SignIn)

		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{3})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": "12345678abc"})
		req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 400, rec.Code)
		assert.NotEmpty(suite.T(), res.Error)

		suite.TearDownTest()
	})

	suite.Run("abnormal.different_password", func() {
		suite.SetupTest()
		suite.e.POST("/user/signin", suite.h.SignIn)

		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{1})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": "12345678abcd"})
		req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 400, rec.Code)
		assert.NotEmpty(suite.T(), res.Error)

		suite.TearDownTest()
	})

	suite.Run("abnormal.different_mail", func() {
		suite.SetupTest()
		suite.e.POST("/user/signin", suite.h.SignIn)

		var expected_user model.User
		suite.seeder.DB.Find(&expected_user, []int64{1})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": "dummy" + expected_user.Mail, "password": "12345678abc"})
		req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		suite.e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(suite.T(), 400, rec.Code)
		assert.NotEmpty(suite.T(), res.Error)

		suite.TearDownTest()
	})

	suite.Run("abnormal.validation", func() {
		suite.SetupTest()
		suite.e.POST("/user/signin", suite.h.SignIn)

		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "", "password": ""})
			req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Len(suite.T(), res.Error.Messages, 2)
		}
		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "aaaaaaa", "password": "12345678abc"})
			req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Len(suite.T(), res.Error.Messages, 1)
		}
		{
			post_data, _ := json.Marshal(map[string]interface{}{"password": "12345678abc"})
			req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Len(suite.T(), res.Error.Messages, 1)

		}

		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "user1@test.com", "password": "12"})
			req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Len(suite.T(), res.Error.Messages, 1)

		}

		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "user1@test.com", "password": "1234567890123456789012345678901"})
			req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Len(suite.T(), res.Error.Messages, 1)

		}

		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "user1@test.com"})
			req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			suite.e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(suite.T(), 400, rec.Code)
			assert.Len(suite.T(), res.Error.Messages, 1)

		}
		suite.TearDownTest()
	})
}
