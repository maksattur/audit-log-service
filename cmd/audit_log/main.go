package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/maksattur/audit-log-service/internal/config"
	"github.com/maksattur/audit-log-service/internal/repository/ch"
	"github.com/maksattur/audit-log-service/internal/services"
	"github.com/maksattur/audit-log-service/internal/token_manager"
	"github.com/maksattur/audit-log-service/internal/transport/consumer"
	"github.com/maksattur/audit-log-service/internal/transport/consumer/kafka_consumer"
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

	eventStoreRepo, err := ch.NewClickHouse(ctx, cfg.ClickHouse)
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse %w", err)
	}

	auditService := services.NewEventService(eventStoreRepo)

	tokenManager, err := token_manager.NewTokenManager(cfg.SecretKey, cfg.JwtTTL)
	if err != nil {
		return fmt.Errorf("failed to make token manager %w", err)
	}

	httpHandler := handler.NewHttpHandler(auditService, tokenManager)

	kafkaCfg := kafka_consumer.NewConfig(cfg.KafkaBrokerAddr, cfg.KafkaGroupID, cfg.KafkaTopic)

	kafkaConsumer := kafka_consumer.NewKafkaConsumer(kafkaCfg)

	consumerAdapter := consumer.NewConsumer(auditService, kafkaConsumer)
	defer consumerAdapter.Close()

	// run receive data
	go consumerAdapter.Receive(ctx)

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
