package main

import (
	"api/config"
	"api/controllers"
	"api/routes"
	"api/setup"
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	secrets := config.GetSecrets()
	logger := logrus.New()

	// Choose the appropriate initializer based on the selected database in the configuration.
	var initializer setup.ServiceInitializer
	switch secrets.CurrentDatabase {
	case "mongodb":
		initializer = setup.MongoDBInitializer{Secrets: secrets, Logger: logger}
	default:
		log.Fatal("Invalid database type specified in configuration.")
	}
	opts, err := setup.ConfigureServiceDependencies(initializer)
	if err != nil {
		log.Fatal("Error configuring service dependencies:", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(opts.Middlewares.AuthMiddleware)

	controller := controllers.NewController(opts)
	routes.ConfigureRoutes(router, controller)

	address := "0.0.0.0:" + secrets.Port
	server := http.Server{
		Addr:    address,
		Handler: router,
	}

	// Handle graceful shutdown on receiving signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Connect to http://localhost:%s/ for server API", secrets.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for a signal to gracefully shut down the server.
	<-stop
	log.Println("Shutting down server...")

	// Create a context with a timeout to force shutdown after a certain duration.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server gracefully stopped")
}
