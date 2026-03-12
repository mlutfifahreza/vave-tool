package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/vave-tool/docs"
	"github.com/vave-tool/internal/api/handler"
	"github.com/vave-tool/internal/api/router"
	"github.com/vave-tool/internal/config"
	grpcHandler "github.com/vave-tool/internal/grpc"
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
	cfg := config.Load()

	database, err := db.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Database connection established")

	redisClient, err := db.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("Redis connection established")

	productRepo := repository.NewProductRepository(database)
	productService := service.NewProductService(productRepo, redisClient)
	productHandler := handler.NewProductHandler(productService)

	httpRouter := router.NewRouter(productHandler)
	httpMux := httpRouter.SetupRoutes()

	grpcServer := grpc.NewServer()
	productGRPCServer := grpcHandler.NewProductServer(productService)
	proto.RegisterProductServiceServer(grpcServer, productGRPCServer)
	reflection.Register(grpcServer)

	go func() {
		httpAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting HTTP server on %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, httpMux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	go func() {
		grpcAddr := fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcAddr, err)
		}
		log.Printf("Starting gRPC server on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	grpcServer.GracefulStop()
	log.Println("Servers stopped")
}
