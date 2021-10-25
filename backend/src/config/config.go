package config

import (
	"embed"
	"encoding/json"
	"os"
)

//go:embed config.prod.json
//go:embed config.stg.json
//go:embed config.dev.json
//go:embed config.test.json
var raw embed.FS

var config *Config

type Config struct {
	Mode    string  `json:"mode"`
	App     App     `json:"app"`
	Log     log     `json:"log"`
	Salt    string  `json:"salt"`
	Web     Web     `json:"web"`
	DB      DB      `json:"db"`
	Session Session `json:"session"`
	Mail    Mail    `json:"mail"`
}

type App struct {
	Name              string            `json:"name"`
	Domain            string            `json:"domain"`
	FrontendDomain    string            `json:"frontend_domain"`
	PreSignUp         preSignUp         `json:"pre_signup"`
	PreForgotPassword preForgotPassword `json:"pre_forgot_password"`
}
type preSignUp struct {
	Lifetime int `json:"lifetime"`
}
type preForgotPassword struct {
	Lifetime int `json:"lifetime"`
}

type log struct {
	DBSlowThreshold int `json:"db_slow_threshold_sec"`
}

type Web struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Cors cors   `json:"cors"`
}
type cors struct {
	AllowOrigins []string `json:"allow_origins"`
}

type DB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db_name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Session struct {
	Codec    string `json:"codec"`
	Key      string `json:"key"`
	Lifetime int    `json:"lifetime"`
	Path     string `json:"path"`
	Table    string `json:"table"`
}

type Mail struct {
	Smtp smtp   `json:"smtp"`
	From string `json:"from"`
}
type smtp struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func init() {
	config = &Config{}
	mode := os.Getenv("MODE")
	b := make([]byte, 0)

	switch mode {
	case "prod":
		b, _ = raw.ReadFile("config.prod.json")
	case "dev":
		b, _ = raw.ReadFile("config.dev.json")
	default:
		b, _ = raw.ReadFile("config.test.json")
	}

	json.Unmarshal([]byte(b), config)
}

func Get() *Config {
	return config
}
