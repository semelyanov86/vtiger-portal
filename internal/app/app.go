package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	http2 "github.com/semelyanov86/vtiger-portal/internal/delivery/http"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/internal/server"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/email/smtp"
	"github.com/semelyanov86/vtiger-portal/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Run initializes whole application.
func Run(configPath string) {
	var wg sync.WaitGroup

	cfg := config.Init(configPath)
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(logger.ConvertErrorToStruct(err, 0, nil))
		return
	}
	memcache := cache.NewMemoryCache()
	emailSender := smtp.NewMailer(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender)

	repos := repository.NewRepositories(db, *cfg, memcache)
	services := service.NewServices(*repos, emailSender, &wg)
	handlers := http2.NewHandler(services, cfg)
	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init())

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error(logger.GenerateErrorMessageFromString("error occurred while running http server:" + err.Error()))
		}
	}()

	logger.Info(logger.GenerateErrorMessageFromString("Server started at port: " + strconv.Itoa(cfg.HTTP.Port)))

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	wg.Wait()

	if err := srv.Stop(ctx); err != nil {
		logger.Error(logger.GenerateErrorMessageFromString("Failed to stop server: " + err.Error()))
	}

	if err := db.Close(); err != nil {
		logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
	}

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
