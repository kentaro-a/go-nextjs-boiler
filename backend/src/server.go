package main

import (
	"app/config"
	app_log "app/log"
	"app/router"
	"fmt"
	"log"
	"runtime"
)

func main() {
	c := config.Get()
	log.Println("Mode: ", c.Mode)
	e, err := router.New()
	if err != nil {
		app_log.Fatal(nil, app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Error:      err,
			Messages:   []string{"Router作成時にエラーが発生しました"},
		})
	}

	err = e.Start(fmt.Sprintf("%s:%d", c.Web.Host, c.Web.Port))
	if err != nil {
		app_log.Fatal(nil, app_log.Fields{
			ScriptInfo: app_log.GetScriptInfo(runtime.Caller(0)),
			Error:      err,
			Messages:   []string{"サーバー開始時にエラーが発生しました"},
		})
	}

}
