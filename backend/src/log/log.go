package log

import (
	"app/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/imdario/mergo"
	echo "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ScriptInfo struct {
	Function string
	File     string
	Line     int
}

func GetScriptInfo(pc uintptr, file string, line int, ok bool) *ScriptInfo {
	if ok {
		script_info := &ScriptInfo{}
		script_info.Function = runtime.FuncForPC(pc).Name()
		script_info.File = file
		script_info.Line = line
		return script_info
	} else {

		return nil
	}
}

type Fields struct {
	File       string      `json:"file"`
	Line       int         `json:"line"`
	Messages   []string    `json:"message"`
	Dump       interface{} `json:"dump"`
	FmtError   []string    `json:"error"`
	ScriptInfo *ScriptInfo `json:"-"`
	Error      error       `json:"-"`
}

func (f Fields) JSON() map[string]interface{} {
	if f.Error != nil {
		errm := strings.Replace(fmt.Sprintf("%+v", f.Error), "\n\t", ": ", -1)
		ferr := strings.Split(errm, "\n")
		f.FmtError = ferr
	}

	if f.ScriptInfo != nil {
		f.Line = f.ScriptInfo.Line
		f.File = f.ScriptInfo.File
	}

	j, _ := json.Marshal(f)
	var m map[string]interface{}
	json.Unmarshal(j, &m)
	return m
}

var cf *config.Config

func init() {
	cf = config.Get()
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "created_at",
			logrus.FieldKeyLevel: "log_level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	logrus.SetOutput(os.Stdout)
	if cf.Mode == "prod" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if cf.Mode == "test" {
		logrus.SetLevel(logrus.DebugLevel)
		// logrus.SetLevel(logrus.PanicLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func getEntry(c echo.Context, fields Fields) *logrus.Entry {
	req_body := []byte{}
	if c != nil {
		if c.Request().Body != nil {
			req_body, _ = ioutil.ReadAll(c.Request().Body)
		}
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(req_body))
		def := map[string]interface{}{
			"mode":         cf.Mode,
			"ip":           c.Request().RemoteAddr,
			"method":       c.Request().Method,
			"path":         c.Request().URL.Path,
			"query":        c.Request().URL.RawQuery,
			"request_body": string(req_body),
			"user_id":      "",
		}
		mergo.MergeWithOverwrite(&def, fields.JSON())
		return logrus.WithFields(def)

	} else {
		def := map[string]interface{}{
			"mode":         cf.Mode,
			"ip":           "",
			"method":       "",
			"path":         "",
			"query":        "",
			"request_body": "",
			"user_id":      "",
		}
		mergo.MergeWithOverwrite(&def, fields.JSON())
		return logrus.WithFields(def)
	}

}

func Info(c echo.Context, fields Fields) {
	getEntry(c, fields).Info()
}

func Warn(c echo.Context, fields Fields) {
	getEntry(c, fields).Warn()
}

func Error(c echo.Context, fields Fields) {
	getEntry(c, fields).Error()
}

func Fatal(c echo.Context, fields Fields) {
	getEntry(c, fields).Fatal()
}

func Panic(c echo.Context, fields Fields) {
	getEntry(c, fields).Panic()
}
