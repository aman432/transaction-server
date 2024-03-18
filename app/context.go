package app

import (
	"context"
	"sync"
	"transaction-server/internal/common/db"
	"transaction-server/internal/config"
)

// this file consists of the functionalities required throughout the app:
// db, logger, producer, tracer
var (
	appContext *applicationContext
	once       sync.Once
)

type applicationContext struct {
	ctx context.Context
	// Application Configuration
	config config.AppConfig

	// DB holds the db connection.
	dB *db.DB
}

// NewAppContext creates New Application Context(singleton function)
func NewAppContext(ctx context.Context, config config.AppConfig) *applicationContext {
	once.Do(func() {
		appContext = new(applicationContext)
		appContext.ctx = ctx
		appContext.config = config
	})

	return appContext
}

// Context To retrieve Application Context
func Context() *applicationContext {
	return appContext
}

func (appContext *applicationContext) Ctx() context.Context {
	return appContext.ctx
}

func (appContext *applicationContext) Config() config.AppConfig {
	return appContext.config
}

func (appContext *applicationContext) DB() *db.DB {
	return appContext.dB
}

func (appContext *applicationContext) SetDB(dB *db.DB) {
	if appContext.dB == nil {
		appContext.dB = dB
	}
}
