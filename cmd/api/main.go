package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/vave-tool/docs"
	"github.com/vave-tool/internal/api/handler"
	"github.com/vave-tool/internal/api/router"
	"github.com/vave-tool/internal/config"
	grpcHandler "github.com/vave-tool/internal/grpc"
	"github.com/vave-tool/internal/observability"
	"github.com/vave-tool/internal/pkg/db"
	"github.com/vave-tool/internal/repository"
	"github.com/vave-tool/internal/service"
	"github.com/vave-tool/proto"
)

// @title Vave Tool API
// @version 1.0
// @description A Go-based backend service that provides REST and gRPC APIs for managing products.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@vave-tool.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "vave-tool-api"
	}

	serviceVersion := os.Getenv("SERVICE_VERSION")
	if serviceVersion == "" {
		serviceVersion = "1.0.0"
	}

	otelEndpoint := os.Getenv("OTEL_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "localhost:4319"
	}

	telemetry, err := observability.InitTelemetry(
		serviceName,
		serviceVersion,
		otelEndpoint,
		cfg.LogLevel,
	)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer func() {
		if err := telemetry.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down telemetry: %v", err)
		}
	}()

	obsLogger := observability.NewLogger(telemetry.Logger)
	telemetry.Logger.Info("Starting Vave Tool API")

	database, err := db.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	telemetry.Logger.Info("Database connection established")

	redisClient, err := db.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	telemetry.Logger.Info("Redis connection established")

	productRepo := repository.NewProductRepository(database)
	productService := service.NewProductService(productRepo, redisClient, obsLogger)
	productHandler := handler.NewProductHandler(productService, obsLogger)

	middleware, err := observability.NewMiddleware(telemetry.Logger)
	if err != nil {
		log.Fatalf("Failed to create middleware: %v", err)
	}

	httpRouter := router.NewRouter(productHandler, middleware, telemetry.MetricsHandler)
	httpMux := httpRouter.SetupRoutes()

	grpcServer := grpc.NewServer()
	productGRPCServer := grpcHandler.NewProductServer(productService)
	proto.RegisterProductServiceServer(grpcServer, productGRPCServer)
	reflection.Register(grpcServer)

	httpAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	httpServer := &http.Server{
		Addr:           httpAddr,
		Handler:        httpMux,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	go func() {
		telemetry.Logger.Info(fmt.Sprintf("Starting HTTP server on %s", httpAddr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	go func() {
		grpcAddr := fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
		}
		telemetry.Logger.Info(fmt.Sprintf("Starting gRPC server on %s", grpcAddr))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	telemetry.Logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		telemetry.Logger.Error(fmt.Sprintf("HTTP server shutdown error: %v", err))
	}
	grpcServer.GracefulStop()
	telemetry.Logger.Info("Servers stopped")
}
