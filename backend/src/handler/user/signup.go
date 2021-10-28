package user

import (
	"app/config"
	app_log "app/log"
	"app/mail"
	"app/model"
	"app/response"
	app_validator "app/validator"
	"fmt"
	"runtime"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) SignUpVerifyToken(c echo.Context) error {
	// middlewareでトークン検証しているので、このハンドラに渡った時点でトークンは正しいことが担保されている
	user_mail_auth := c.Get("user_mail_auth").(model.MailAuth)
	return response.Success(c, 200, map[string]interface{}{
		"mail":      user_mail_auth.Mail,
		"function":  user_mail_auth.Function,
		"token":     user_mail_auth.Token,
		"expire_at": user_mail_auth.ExpireAt,
	}, nil)

}

func (h Handler) SignUp(c echo.Context) error {
	user_mail_auth := c.Get("user_mail_auth").(model.MailAuth)
	user_model := model.NewUserModel(h.DB)
	user := model.User{}
	c.Bind(&user)
	user.Mail = user_mail_auth.Mail
	user.StatusFlg = 0

	vld := app_validator.Get()
	err := vld.Struct(user)
	if err != nil {
		messages := app_validator.GetErrorMessages(&model.User{}, err)
		return response.Error(c, 400, messages, nil)
	}

	is_exist, _, err := user_model.IsMailExist(user.Mail)
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
	user_mail_auth_model := model.NewMailAuthModel(tx)
	user_model = model.NewUserModel(tx)

	// usersに登録
	user.Password = user_model.GetHashedPassword(user.Password)
	err = user_model.Create(&user)
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
	sender := mail.NewSender("signup", user.Mail, "会員登録完了のお知らせ", map[string]string{
		"@NAME@":       user.Name,
		"@SIGNIN_URL@": fmt.Sprintf("%ssignin", config.Get().App.Domain),
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
