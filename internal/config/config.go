package config

import (
	"transaction-server/internal/common/db"
)

type AppConfig struct {
	App App
	Db  db.Config
}

type App struct {
	ServiceName string
	Hostname    string
	Port        string
}
