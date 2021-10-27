package handler

import (
	"app/config"
	app_log "app/log"
	"app/mail"
	"app/model"
	"app/response"
	"app/util"
	"fmt"
	"runtime"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) ForgotPasswordVerifyToken(c echo.Context) error {
	// middlewareでトークン検証しているので、このハンドラに渡った時点でトークンは正しいことが担保されている
	user_mail_auth := c.Get("user_mail_auth").(model.UserMailAuth)
	return response.Success(c, 200, map[string]interface{}{
		"mail":      user_mail_auth.Mail,
		"function":  user_mail_auth.Function,
		"token":     user_mail_auth.Token,
		"expire_at": user_mail_auth.ExpireAt,
	}, nil)

}

func (h Handler) ForgotPassword(c echo.Context) error {
	user_mail_auth := c.Get("user_mail_auth").(model.UserMailAuth)
	user_model := model.NewUserModel(h.DB)

	is_exist, user, err := user_model.IsMailExist(user_mail_auth.Mail)
	if err != nil {
		return response.SystemError(c, &app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Messages:   []string{err.Error()},
			Error:      err,
		})
	}
	if !is_exist {
		return response.Error(c, 400, []string{"メールアドレスが登録されていません"}, nil)
	}

	tx := h.DB.Begin()
	user_mail_auth_model := model.NewUserMailAuthModel(tx)
	user_model = model.NewUserModel(tx)

	// ランダムパスワード生成
	new_password := util.MakeRandStr(20)

	// パスワード更新
	user.Password = user_model.GetHashedPassword(new_password)
	err = user_model.Save(user)
	if err != nil {
		tx.Rollback()
		return response.SystemError(c, &app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Messages:   []string{err.Error()},
			Error:      err,
		})
	}

	// すでに登録されてる仮登録データを削除
	err = user_mail_auth_model.DeleteByMailFunction(user_mail_auth.Mail, user_mail_auth.Function)
	if err != nil {
		tx.Rollback()
		return response.SystemError(c, &app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Messages:   []string{err.Error()},
			Error:      err,
		})
	}

	// トークン付きの本登録URLをメールで送信
	sender := mail.NewSender("forgot_password", user.Mail, "パスワード再発行完了のお知らせ", map[string]string{
		"@NAME@":         user.Name,
		"@SIGNIN_URL@":   fmt.Sprintf("%ssignin", config.Get().App.Domain),
		"@NEW_PASSWORD@": new_password,
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
