package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"gitlab.com/InfoBlogFriends/server/components"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/logs"

	"github.com/subosito/gotenv"

	"gitlab.com/InfoBlogFriends/server/db"

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

	cmps := components.NewComponents(logger)

	repositories := db.NewRepositories(cmps)

	service := services.NewServices(cmps, repositories)

	// router initialization
	r, err := server.NewRouter(service, cmps)

	if err != nil {
		logger.Fatal("router initialization error", zap.Error(err))
	}

	// server initialization
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cmps.Config().Server.Port),
		Handler: r,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http listen and serve error", zap.Error(err))
		}
	}()

	logger.Info("server started", zap.Int("port", cmps.Config().Server.Port))

	<-ctx.Done()

	logger.Info("server stopped")
}
