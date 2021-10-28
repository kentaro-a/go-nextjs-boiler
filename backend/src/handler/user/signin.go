package user

import (
	"app/config"
	app_log "app/log"
	"app/model"
	"app/response"
	"app/session"
	app_validator "app/validator"
	"encoding/json"
	"runtime"

	"github.com/pkg/errors"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h Handler) SignIn(c echo.Context) error {

	conf := config.Get()
	m := model.NewUserModel(h.DB)
	posted_data := model.User{}
	c.Bind(&posted_data)

	vld := app_validator.Get()
	err := vld.StructPartial(posted_data, []string{"Mail", "Password"}...)
	if err != nil {
		messages := app_validator.GetErrorMessages(&model.User{}, err)
		return response.Error(c, 400, messages, nil)
	}

	// status_flg = 0などstructで検索するとゼロ値は検索条件から除外されるので直接指定する
	user, err := m.FindByMailPassword(
		posted_data.Mail,
		m.GetHashedPassword(posted_data.Password),
	)

	if err == nil {
		sess, _ := session.Get(c, conf.Session.Key)
		user_without_secrets := user.NewUserWithoutSecrets()
		buser, _ := json.Marshal(user_without_secrets)
		sess.Values["signin_user"] = buser
		sess.Save(c.Request(), c.Response())

	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return response.SystemError(c, &app_log.Fields{
				ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
				Messages:   []string{err.Error()},
				Error:      err,
			})
		} else {
			return response.Error(c, 400, []string{"メールアドレスまたはパスワードが違います"}, nil)
		}
	}

	return response.Success(c, 200, nil, nil)
}
