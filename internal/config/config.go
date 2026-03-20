package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	GRPC     GRPCConfig
	Auth     AuthConfig
	LogLevel string
}

type ServerConfig struct {
	Port           string
	Host           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	Host            string
	Port            string
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	PoolTimeout     time.Duration
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type GRPCConfig struct {
	Port string
	Host string
}

type AuthConfig struct {
	InternalUsername string
	InternalPassword string
	JWTSecret        string
	GoogleClientID   string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           getEnv("SERVER_PORT", "8080"),
			Host:           getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:    getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:   getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:    getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
			MaxHeaderBytes: getIntEnv("SERVER_MAX_HEADER_BYTES", 1048576),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "vave_db"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:            getEnv("REDIS_HOST", "localhost"),
			Port:            getEnv("REDIS_PORT", "6379"),
			Password:        getEnv("REDIS_PASSWORD", ""),
			DB:              0,
			PoolSize:        getIntEnv("REDIS_POOL_SIZE", 10),
			MinIdleConns:    getIntEnv("REDIS_MIN_IDLE_CONNS", 2),
			PoolTimeout:     getDurationEnv("REDIS_POOL_TIMEOUT", 4*time.Second),
			ConnMaxLifetime: getDurationEnv("REDIS_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDurationEnv("REDIS_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50051"),
			Host: getEnv("GRPC_HOST", "0.0.0.0"),
		},
		Auth: AuthConfig{
			InternalUsername: getEnv("INTERNAL_AUTH_USERNAME", "admin"),
			InternalPassword: getEnv("INTERNAL_AUTH_PASSWORD", "admin"),
			JWTSecret:        getEnv("JWT_SECRET", ""),
			GoogleClientID:   getEnv("GOOGLE_CLIENT_ID", ""),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
