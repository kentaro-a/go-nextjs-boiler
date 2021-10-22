package handler

import (
	"app/model"
	"app/response"
	"app/util"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	e, h, _, seeder := setup(t)

	e.POST("/signin", h.SignIn)

	// 正常
	{
		var expected_user model.User
		seeder.DB.Find(&expected_user, []int64{1})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": "12345678abc"})
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, res.StatusCode, rec.Code)
		assert.NotEmpty(t, res.Data)

		d := res.Data.(map[string]interface{})
		var actual_user model.User
		util.MapToStruct(d["user"], &actual_user)
		assert.Equal(t, expected_user.ID, actual_user.ID)
		assert.NotEmpty(t, rec.Header().Get("Set-Cookie"))
	}
	// status_flg=0のユーザー
	{
		var expected_user model.User
		seeder.DB.Find(&expected_user, []int64{3})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": "12345678abc"})
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 400, rec.Code)
		assert.NotEmpty(t, res.Error)
	}
	// password違い
	{
		var expected_user model.User
		seeder.DB.Find(&expected_user, []int64{1})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": "12345678abcd"})
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 400, rec.Code)
		assert.NotEmpty(t, res.Error)

	}

	// mail違い
	{
		var expected_user model.User
		seeder.DB.Find(&expected_user, []int64{1})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": "dummy" + expected_user.Mail, "password": "12345678abc"})
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseError{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 400, rec.Code)
		assert.NotEmpty(t, res.Error)

	}

	// validation
	{

		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "", "password": ""})
			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Len(t, res.Error.Messages, 2)
		}
		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "aaaaaaa", "password": "12345678abc"})
			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Len(t, res.Error.Messages, 1)
		}
		{
			post_data, _ := json.Marshal(map[string]interface{}{"password": "12345678abc"})
			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Len(t, res.Error.Messages, 1)

		}

		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "user1@test.com", "password": "12"})
			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Len(t, res.Error.Messages, 1)

		}

		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "user1@test.com", "password": "1234567890123456789012345678901"})
			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Len(t, res.Error.Messages, 1)

		}

		{
			post_data, _ := json.Marshal(map[string]interface{}{"mail": "user1@test.com"})
			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			res := response.ResponseError{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, 400, rec.Code)
			assert.Len(t, res.Error.Messages, 1)

		}
	}
	teardown(t, e, seeder)
}
