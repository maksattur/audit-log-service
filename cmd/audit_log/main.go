package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/maksattur/audit-log-service/internal/config"
	"github.com/maksattur/audit-log-service/internal/transport/handler"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// listen to OS signals and gracefully shutdown application
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load config %w", err)
	}

	httpHandler := handler.NewHttpHandler()

	srv := http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpHandler.Router(),
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
	}()

	log.Printf("HTTP server is starting on %s", cfg.HTTPAddr)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	log.Printf("Have a nice day!")

	return nil
}
