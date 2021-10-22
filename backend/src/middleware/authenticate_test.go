package middleware

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

func TestUserAuthentication(t *testing.T) {
	e, h, m, seeder := setup(t)

	e.POST("/signin", h.SignIn)
	e.GET("/dashboard", h.Dashboard, m.UserAuthentication)

	// Has session
	{
		var expected_user model.User
		seeder.DB.Find(&expected_user, []int64{1})
		post_data, _ := json.Marshal(map[string]interface{}{"mail": expected_user.Mail, "password": "12345678abc"})
		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(post_data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		cookie := rec.Header().Get("Set-Cookie")

		req = httptest.NewRequest(http.MethodGet, "/dashboard", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Cookie", cookie)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)

	}

	// Invalid session
	{

		cookie := "mapp=MTYzMzA3MDQ0MHxCQXdBQVRJPXwReJY1tkoCaIDBycrcKTa8n3tJHMKieKhuuOjrLrpDiQ==; Path=/; Expires=Fri, 01 Oct 2021 07:40:40 GMT; Max-Age=3600"

		req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Cookie", cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 401, rec.Code)

	}

	// Missed cookie
	{

		req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 401, rec.Code)

	}

	teardown(t, e, seeder)
}
