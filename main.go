package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/server"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), common.Timeouts.Server.Startup)
	defer cancel()

	s, err := server.NewServer(ctx)
	if err != nil {
		slog.Error("failed to create server", "error", err)
		os.Exit(1)
	}

	srvErrsCh := make(chan error, 1)

	go func() {
		srvErrsCh <- s.Start()
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-srvErrsCh:
		if !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Error("server error", "error", err)
		}
	case <-shutdownCh:
		ctx, cancel := context.WithTimeout(context.Background(), common.Timeouts.Server.Write)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			s.Logger.Error("graceful shutdown failed", "error", err)
		}
	}

	s.Logger.Info("Server stopped gracefully")
}
