package user

import (
	"app/config"
	app_context "app/context"
	app_log "app/log"
	"app/mail"
	"app/model"
	"app/response"
	"app/session"
	"runtime"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	echo "github.com/labstack/echo/v4"
)

func (h Handler) Delete(c echo.Context) error {
	cc := app_context.CastContext(c)

	tx := h.DB.Begin()
	user_model := model.NewUserModel(tx)

	posted_data := model.User{}
	c.Bind(&posted_data)

	password := posted_data.Password
	user, err := user_model.FindByID(cc.User.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return response.SystemError(c, &app_log.Fields{
				ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
				Messages:   []string{err.Error()},
				Error:      err,
			})
		} else {
			return response.Error(c, 400, []string{"ユーザーが存在しません"}, nil)
		}
	}

	// パスワード認証
	if user.Password != user_model.GetHashedPassword(password) {
		return response.Error(c, 400, []string{"パスワードが一致しません"}, nil)
	}

	err = user_model.Delete(user)
	if err != nil {
		return response.SystemError(c, &app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Messages:   []string{err.Error()},
			Error:      err,
		})
	}

	// トークン付きの本登録URLをメールで送信
	sender := mail.NewSender("user/delete", user.Mail, "退会完了のお知らせ", map[string]string{
		"@NAME@": user.Name,
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

	tx.Commit()

	// TODO: Not work
	// セッション削除
	session.DeleteSession(c, config.Get().Session.Key)

	return response.Success(c, 200, nil, nil)
}
