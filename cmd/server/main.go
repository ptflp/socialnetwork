package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"gitlab.com/ptflp/infoblog-server/cache"

	"gitlab.com/ptflp/infoblog-server/auth"
	"gitlab.com/ptflp/infoblog-server/db"

	"gitlab.com/ptflp/infoblog-server/respond"

	"gitlab.com/ptflp/infoblog-server/config"
	"gitlab.com/ptflp/infoblog-server/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	GeneralError = 2
)

func main() {
	// config initialization
	zapConf := zap.NewProductionConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConf.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	atom := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zapConf.Level = atom
	logger, err := zapConf.Build()
	if err != nil {
		fmt.Printf("logger initialization: %s\n", err)
		os.Exit(GeneralError)
	}

	logger = logger.Named("infoBlog")

	// shutdown server on signal interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigInt := <-sig
		logger.Info("signal interrupt recieved", zap.Stringer("os signal", sigInt))
		cancel()
	}()

	responder, err := respond.NewResponder(logger)
	if err != nil {
		logger.Fatal("responder initialization error", zap.Error(err))
	}

	// config initialization
	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatal("config initialization error", zap.Error(err))
	}

	c, err := cache.NewRedisCache(conf.Redis)

	database, err := db.NewDB(logger, conf.DB)

	userRepository := db.NewUserRepository(database)
	authService := auth.NewAuthService(userRepository, c)

	// router initialization
	r, err := server.NewRouter(&server.Services{AuthService: authService}, &server.HandlerComponents{
		UserRepository: nil,
		Logger:         logger,
		Responder:      responder,
		LogLevel:       zap.AtomicLevel{},
	})
	if err != nil {
		logger.Fatal("router initialization error", zap.Error(err))
	}

	// server initialization
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: r,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http listen and serve error", zap.Error(err))
		}
	}()

	logger.Info("server started", zap.Int("port", conf.Server.Port))

	<-ctx.Done()

	logger.Info("server stopped")
}
