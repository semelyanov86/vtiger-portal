package app

import (
	"context"
	"database/sql"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/logger"
	"time"
)

// @title Customer Portal for Vtiger
// @version 1.0
// @description REST API for Customer Portal

// @host localhost:4050
// @BasePath /api/v1/

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes whole application.
func Run(configPath string) {
	cfg := config.Init(configPath)
	_, err := openDB(cfg)
	if err != nil {
		logger.Error(logger.ConvertErrorToStruct(err, 0, nil))
		return
	}
	_ = cache.NewMemoryCache()
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
