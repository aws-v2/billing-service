package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Qarani-m/billing-service/internal/application"
	"github.com/Qarani-m/billing-service/internal/config"
	"github.com/Qarani-m/billing-service/internal/infrastructure/messaging"
	"github.com/Qarani-m/billing-service/internal/infrastructure/repository"
	transport "github.com/Qarani-m/billing-service/internal/transport/http"
	"github.com/Qarani-m/billing-service/pkg/database"
	internalMessaging "github.com/Qarani-m/billing-service/pkg/messaging"
	"github.com/gorilla/mux"
)

func main() {
	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Setup Database
	dbConfig := database.PostgresConfig{
		Host:            cfg.DB.Host,
		Port:            cfg.DB.Port,
		User:            cfg.DB.User,
		Password:        cfg.DB.Password,
		Database:        cfg.DB.Database,
		SSLMode:         cfg.DB.SSLMode,
		MaxOpenConns:    cfg.DB.MaxOpenConns,
		MaxIdleConns:    cfg.DB.MaxIdleConns,
		ConnMaxLifetime: cfg.DB.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.DB.ConnMaxIdleTime,
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 3. Setup NATS
	natsPublisher, err := internalMessaging.NewNATSPublisher(
		cfg.NATS.URL,
		cfg.NATS.User,
		cfg.NATS.Password,
		cfg.Profile,
	)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsPublisher.Close()

	// 4. Initialize Dependency Injection
	repo := repository.NewPostgresInvoiceRepository(db)
	eventPublisher := messaging.NewNATSEventPublisher(natsPublisher)
	service := application.NewInvoiceService(repo, eventPublisher)
	handler := transport.NewHandler(service)

	// 5. Setup Router
	r := mux.NewRouter()
	handler.RegisterRoutes(r)

	// Add a simple health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 6. Start Server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting Billing Service on port %s in %s mode", cfg.Server.Port, cfg.Profile)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
