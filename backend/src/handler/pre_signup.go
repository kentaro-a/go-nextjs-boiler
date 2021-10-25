package handler

import (
	"app/config"
	app_log "app/log"
	"app/mail"
	"app/model"
	"app/response"
	"app/util"
	app_validator "app/validator"
	"fmt"
	"runtime"
	"time"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) PreSignUp(c echo.Context) error {

	user_model := model.NewUserModel(h.DB)
	user := model.User{}
	c.Bind(&user)

	vld := app_validator.Get()
	err := vld.StructPartial(user, "Mail")
	if err != nil {
		messages := app_validator.GetErrorMessages(&model.User{}, err)
		return response.Error(c, 400, messages, nil)
	}

	is_exist, err := user_model.IsMailExist(user.Mail)
	if err != nil {
		return response.SystemError(c, &app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Messages:   []string{err.Error()},
			Error:      err,
		})
	}
	if is_exist {
		return response.Error(c, 400, []string{"すでに登録済みのメールアドレスです"}, nil)
	}

	tx := h.DB.Begin()
	user_mail_auth_model := model.NewUserMailAuthModel(tx)
	user_mail_auth := model.UserMailAuth{
		Function: "pre_signup",
		Mail:     user.Mail,
		Token:    util.MakeRandStr(62),
		ExpireAt: time.Now().Add(time.Second * time.Duration(config.Get().App.PreSignUp.Lifetime)), // 有効期限
	}

	// まずすでに登録されてる仮登録データを削除
	err = user_mail_auth_model.DeleteByMailFunction(user_mail_auth.Mail, user_mail_auth.Function)
	if err != nil {
		tx.Rollback()
		return response.SystemError(c, &app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Messages:   []string{err.Error()},
			Error:      err,
		})
	}

	// 仮登録データを登録
	err = user_mail_auth_model.Create(&user_mail_auth)
	if err != nil {
		tx.Rollback()
		return response.SystemError(c, &app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Messages:   []string{err.Error()},
			Error:      err,
		})
	}

	// トークン付きの本登録URLをメールで送信
	sender := mail.NewSender("pre_signup", user_mail_auth.Mail, "仮登録完了のお知らせ", map[string]string{
		"@MAIL@":       user_mail_auth.Mail,
		"@SIGNUP_URL@": fmt.Sprintf("%ssignup/%s", config.Get().App.FrontendDomain, user_mail_auth.Token),
	})
	err = sender.Send()
	if err != nil {
		tx.Rollback()
		return response.SystemError(c, &app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Messages:   []string{err.Error()},
			Error:      err,
		})
	}

	// Commit transaction
	tx.Commit()
	return response.Success(c, 200, nil, nil)
}
