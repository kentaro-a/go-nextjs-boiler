package middleware

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

func TestVerifyMailAuth(t *testing.T) {
	e, _, m, seeder := setup(t)

	e.POST("/test/:token", func(c echo.Context) error {
		return response.Success(c, 200, nil, nil)
	}, m.VerifyMailAuth)

	// valid token
	{
		var expected_user_mail_auth model.UserMailAuth
		seeder.DB.Find(&expected_user_mail_auth, []int64{1})
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/test/%s", expected_user_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 200, rec.Code)

	}
	// expired token
	{
		var expected_user_mail_auth model.UserMailAuth
		seeder.DB.Find(&expected_user_mail_auth, []int64{2})
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/test/%s", expected_user_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 401, rec.Code)

	}
	// inactive token
	{
		var expected_user_mail_auth model.UserMailAuth
		seeder.DB.Find(&expected_user_mail_auth, []int64{3})
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/test/%s", expected_user_mail_auth.Token), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := response.ResponseSuccess{}
		json.NewDecoder(rec.Body).Decode(&res)
		assert.Equal(t, 401, rec.Code)

	}

	teardown(t, e, seeder)
}
