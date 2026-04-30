package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/joho/godotenv"
)

type Config struct {
	DB      DBConfig
	NATS    NATSConfig
	Server  ServerConfig
	Eureka  EurekaConfig
	Profile string
	Daraja  domain.DarajaConfig
}

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type NATSConfig struct {
	URL      string
	User     string
	Password string
}

type ServerConfig struct {
	Port string
}

type EurekaConfig struct {
	ServerURL         string
	AppName           string
	HostName          string
	IPAddr            string
	Port              int
	VipAddress        string
	InstanceID        string
	HeartbeatInterval time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := getEnvInt("PORT", 8061)

	cfg := &Config{
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres-staging-user"),
			Password:        getEnv("DB_PASSWORD", "postgres-staging-password"),
			Database:        getEnv("DB_NAME", "billing_db"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
		},
		NATS: NATSConfig{
			URL:      getEnv("NATS_URL", "nats://localhost:4222"),
			User:     getEnv("NATS_USER", "auth-server"),
			Password: getEnv("NATS_PASSWORD", "auth-secreta"),
		},
		Server: ServerConfig{
			Port: strconv.Itoa(port),
		},
		Eureka: EurekaConfig{
			ServerURL:         getEnv("EUREKA_SERVER_URL", "http://localhost:8761/eureka"),
			AppName:           getEnv("EUREKA_APP_NAME", "BILLING-SERVICE"),
			HostName:          getEnv("EUREKA_HOSTNAME", "localhost"),
			IPAddr:            getEnv("EUREKA_IP_ADDR", "127.0.0.1"),
			Port:              port,
			VipAddress:        getEnv("EUREKA_VIP_ADDRESS", "billing-service"),
			InstanceID:        getEnv("EUREKA_INSTANCE_ID", fmt.Sprintf("billing-service:%d", port)),
			HeartbeatInterval: getEnvDuration("EUREKA_HEARTBEAT_INTERVAL", 30*time.Second),
		},
		Profile: strings.ToLower(getEnv("APP_PROFILE", "dev")),
		Daraja: domain.DarajaConfig{
			ConsumerKey:    getEnv("DARJA_CONSUMER_KEY", ""),
			ConsumerSecret: getEnv("DARJA_CONSUMER_SECRET", ""),
			ShortCode:      getEnv("DARJA_SHORT_CODE", ""),
			PassKey:        getEnv("DARJA_PASS_KEY", ""),
			CallbackURL:    getEnv("DARJA_CALLBACK_URL", ""),
			Environment:    getEnv("DARJA_ENVIRONMENT", "dev"),
			DarajaBaseURL:  getEnv("DARJA_BASE_URL", "http://localhost:8061"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
