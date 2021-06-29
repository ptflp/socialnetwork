package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"gitlab.com/InfoBlogFriends/server/session"

	"gitlab.com/InfoBlogFriends/server/logs"

	"github.com/subosito/gotenv"

	"gitlab.com/InfoBlogFriends/server/cache"

	"gitlab.com/InfoBlogFriends/server/auth"
	"gitlab.com/InfoBlogFriends/server/db"

	"gitlab.com/InfoBlogFriends/server/respond"

	"gitlab.com/InfoBlogFriends/server/config"
	"gitlab.com/InfoBlogFriends/server/server"
	"go.uber.org/zap"
)

func main() {
	// logs initialization
	logger := logs.NewLogger()

	// environment initialization
	err := gotenv.Load()
	if err != nil {
		logger.Fatal("env initialization error", zap.Error(err))
	}

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

	jwt, err := session.NewJWTKeys(logger)

	c, err := cache.NewRedisCache(conf.Redis)
	if err != nil {
		logger.Fatal("redis initialization error", zap.Error(err))
	}

	database, err := db.NewDB(logger, conf.DB)
	if err != nil {
		logger.Fatal("db initialization error", zap.Error(err))
	}

	userRepository := db.NewUserRepository(database)
	authService := auth.NewAuthService(conf.App, userRepository, c, logger, jwt)

	// router initialization
	r, err := server.NewRouter(&server.Services{AuthService: authService}, &server.HandlerComponents{
		UserRepository: nil,
		Logger:         logger,
		Responder:      responder,
		LogLevel:       zap.NewAtomicLevel(),
		JWTKeys:        jwt,
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
