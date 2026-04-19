package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/gin-gonic/gin"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

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

	// 5. Docs
	docsService := application.NewDocsService(getEnv("DOCS_PATH", "./docs"))
	docsHandler := transport.NewDocsHandler(docsService)

	// 6. Setup Gin Router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	handler.RegisterRoutes(router)
	transport.SetupRoutes(router, docsHandler)

	// 7. Start Server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
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

	// 8. Register with Eureka
	for i := 1; i <= 3; i++ {
		if err := registerWithEureka(cfg.Eureka); err != nil {
			log.Printf("Eureka registration attempt %d failed: %v", i, err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	go sendHeartbeat(cfg.Eureka)

	// 9. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := deregisterFromEureka(cfg.Eureka); err != nil {
		log.Printf("Failed to deregister from Eureka: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// ── Eureka ────────────────────────────────────────────────────────────────────

func registerWithEureka(cfg config.EurekaConfig) error {
	payload := map[string]any{
		"instance": map[string]any{
			"instanceId": cfg.InstanceID,
			"hostName":   cfg.HostName,
			"app":        cfg.AppName,
			"ipAddr":     cfg.IPAddr,
			"vipAddress": cfg.VipAddress,
			"status":     "UP",
			"port": map[string]any{
				"$":        cfg.Port,
				"@enabled": "true",
			},
			"dataCenterInfo": map[string]any{
				"@class": "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
				"name":   "MyOwn",
			},
			"healthCheckUrl": fmt.Sprintf("http://%s:%d/health", cfg.HostName, cfg.Port),
			"statusPageUrl":  fmt.Sprintf("http://%s:%d/health", cfg.HostName, cfg.Port),
			"homePageUrl":    fmt.Sprintf("http://%s:%d/", cfg.HostName, cfg.Port),
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	url := fmt.Sprintf("%s/apps/%s", cfg.ServerURL, cfg.AppName)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request build failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("eureka returned %d", resp.StatusCode)
	}

	log.Printf("Registered with Eureka: %s (%s)", cfg.InstanceID, url)
	return nil
}

func sendHeartbeat(cfg config.EurekaConfig) {
	ticker := time.NewTicker(cfg.HeartbeatInterval)
	defer ticker.Stop()

	url := fmt.Sprintf("%s/apps/%s/%s", cfg.ServerURL, cfg.AppName, cfg.InstanceID)
	client := &http.Client{Timeout: 5 * time.Second}

	for range ticker.C {
		req, err := http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			log.Printf("Heartbeat request build failed: %v", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Heartbeat failed: %v", err)
			continue
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			log.Printf("Heartbeat rejected — status %d", resp.StatusCode)
		} else {
			log.Printf("Heartbeat sent: %s", cfg.InstanceID)
		}
		resp.Body.Close()
	}
}

func deregisterFromEureka(cfg config.EurekaConfig) error {
	url := fmt.Sprintf("%s/apps/%s/%s", cfg.ServerURL, cfg.AppName, cfg.InstanceID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("request build failed: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("eureka returned %d", resp.StatusCode)
	}

	log.Printf("Deregistered from Eureka: %s", cfg.InstanceID)
	return nil
}