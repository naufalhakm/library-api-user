package main

import (
	"library-api-user/internal/config"
	"library-api-user/internal/factory"
	"library-api-user/internal/grpc/handlers"
	"library-api-user/internal/routes"
	"library-api-user/pkg/database"
	"library-api-user/proto/auth"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func main() {
	config.LoadConfig()
	psqlDB, err := database.NewPqSQLClient()
	if err != nil {
		log.Fatal("Could not connect to PqSQL:", err)
	}

	provider := factory.InitFactory(psqlDB)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		runGRPCServer()
	}()

	go func() {
		defer wg.Done()
		runHTTPServer(provider)
	}()

	wg.Wait()
}

func runGRPCServer() {
	listener, err := net.Listen("tcp", ":"+config.ENV.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", config.ENV.GRPCPort, err)
	}

	grpcServer := grpc.NewServer()

	authHandler := handlers.NewAuthService()
	auth.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Printf("gRPC server running on port %s\n", config.ENV.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func runHTTPServer(provider *factory.Provider) {
	// Register HTTP routes
	router := routes.RegisterRoutes(provider)

	log.Printf("REST API server running on port %s\n", config.ENV.ServerPort)
	if err := router.Run(":" + config.ENV.ServerPort); err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}
