package tests

import (
	"app/config"
	"app/model"
	"time"
)

func UserMailAuthsFixture() []*model.UserMailAuth {
	return []*model.UserMailAuth{
		{
			ID:        1,
			Function:  "pre_signup",
			Mail:      "pre_signup_user1@test.com",
			Token:     "pre_signup_ValidToken",
			StatusFlg: 0,
			ExpireAt:  time.Now().Add(time.Second * time.Duration(config.Get().App.PreSignUp.Lifetime)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Function:  "pre_signup",
			Mail:      "pre_signup_user2@test.com",
			Token:     "pre_signup_ExpiredToken",
			StatusFlg: 0,
			ExpireAt:  time.Now().Add(time.Hour * 24 * -1).Add(time.Second * time.Duration(config.Get().App.PreSignUp.Lifetime)),
			CreatedAt: time.Now().Add(time.Hour * 24 * -1),
			UpdatedAt: time.Now().Add(time.Hour * 24 * -1),
		},
		{
			ID:        3,
			Function:  "pre_signup",
			Mail:      "pre_signup_user3@test.com",
			Token:     "pre_signup_InActiveToken",
			StatusFlg: 1,
			ExpireAt:  time.Now().Add(time.Second * time.Duration(config.Get().App.PreSignUp.Lifetime)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		{
			ID:        4,
			Function:  "pre_forgot_password",
			Mail:      "pre_forgot_password_user1@test.com",
			Token:     "pre_forgot_password_ValidToken",
			StatusFlg: 0,
			ExpireAt:  time.Now().Add(time.Second * time.Duration(config.Get().App.PreForgotPassword.Lifetime)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        5,
			Function:  "pre_forgot_password",
			Mail:      "pre_forgot_password_user2@test.com",
			Token:     "pre_forgot_password_ExpiredToken",
			StatusFlg: 0,
			ExpireAt:  time.Now().Add(time.Hour * 24 * -1).Add(time.Second * time.Duration(config.Get().App.PreForgotPassword.Lifetime)),
			CreatedAt: time.Now().Add(time.Hour * 24 * -1),
			UpdatedAt: time.Now().Add(time.Hour * 24 * -1),
		},
		{
			ID:        6,
			Function:  "pre_forgot_password",
			Mail:      "pre_forgot_password_user3@test.com",
			Token:     "pre_forgot_password_InActiveToken",
			StatusFlg: 1,
			ExpireAt:  time.Now().Add(time.Second * time.Duration(config.Get().App.PreForgotPassword.Lifetime)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

}
