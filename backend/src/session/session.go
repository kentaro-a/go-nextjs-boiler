package session

import (
	"app/config"
	app_log "app/log"
	"app/model"
	"encoding/json"
	"runtime"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
	"github.com/srinathgs/mysqlstore"
	"gorm.io/gorm"
)

var store *mysqlstore.MySQLStore

func GetStore(db *gorm.DB) *mysqlstore.MySQLStore {
	if store == nil {
		c := config.Get()
		db, err := db.DB()
		if err != nil {
			app_log.Fatal(nil, app_log.Fields{
				ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
				Messages:   []string{"Session Store取得時のDBエラー"},
				Error:      err,
			})

		}
		store, err = mysqlstore.NewMySQLStoreFromConnection(db, c.Session.Table, c.Session.Path, c.Session.Lifetime, []byte(c.Session.Codec))
		if err != nil {
			app_log.Fatal(nil, app_log.Fields{
				ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
				Messages:   []string{"Session Store取得時のDBエラー"},
				Error:      err,
			})
		}
	} else {
	}
	return store
}

func DeleteStore() {
	store = nil
}

func Get(c echo.Context, key string) (*sessions.Session, error) {
	sess, err := session.Get(key, c)
	return sess, err
}

func DeleteSession(c echo.Context, key string) {
	sess, err := session.Get(key, c)
	if err == nil {
		sess.Options.MaxAge = -1
		sess.Save(c.Request(), c.Response())
	}
}

func IsSignedIn(c echo.Context) (bool, model.User) {
	conf := config.Get()
	sess, _ := session.Get(conf.Session.Key, c)
	user := model.User{}
	if sess.Values["signin_user"] == nil {
		return false, user
	} else {
		b, _ := sess.Values["signin_user"].([]byte)
		json.Unmarshal(b, &user)
		return true, user
	}
}
