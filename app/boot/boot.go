package boot

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"transaction-server/app"
	"transaction-server/internal/common/db"
	config2 "transaction-server/internal/config"
)

var (
	// Db holds the application db connection.
	Db *db.DB
	// Config holds info about growth config
	Config config2.AppConfig
)

func init() {
	fmt.Println("Initializing App ...", "default")
	Config = initConfig()
}

func initConfig() config2.AppConfig {
	err := config2.NewDefaultConfig().Load("default", &Config)
	if err != nil {
		log.Fatal(err)
	}
	return Config
}

func Initialize(ctx context.Context) error {
	// This function is used to initialize the application
	var err error

	// Init Db
	Db, err = InitDb()
	if err != nil {
		return err
	}

	// Create an Application appContext which will be used across the application.
	appContext := app.NewAppContext(ctx, Config)
	appContext.SetDB(Db)
	return nil
}

// InitDb initializes db
func InitDb() (*db.DB, error) {
	gDb, err := db.NewDb(&Config.Db, db.GormConfig(getGormConfig(&Config.Db)), db.Dialector(getGormDialector()))
	return gDb, err
}

func getGormConfig(cr db.IConfigReader) *gorm.Config {
	return &gorm.Config{
		AllowGlobalUpdate:        false,
		SkipDefaultTransaction:   true,
		PrepareStmt:              true,
		DisableNestedTransaction: true,
		Logger:                   gormLogger.Default.LogMode(getDbLogLevelByDebugMode(cr.IsDebugMode())),
	}
}

func getDbLogLevelByDebugMode(debug bool) gormLogger.LogLevel {
	if debug == false {
		return gormLogger.Silent
	}
	return gormLogger.Info
}

func getGormDialector() gorm.Dialector {
	dsnFormat := "%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	return mysql.Open(fmt.Sprintf(dsnFormat, Config.Db.Username, Config.Db.Password, Config.Db.Protocol, Config.Db.URL, Config.Db.Port, Config.Db.Name))
}
