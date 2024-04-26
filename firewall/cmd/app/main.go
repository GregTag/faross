package main

import (
	"context"
	"firewall/internal/handler"
	"firewall/internal/service"
	"firewall/internal/storage"
	"firewall/internal/storage/driver"
	"firewall/pkg/config"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load config
	config_path := flag.String("config", "config/config.yaml", "Path to configuration file")
	flag.Parse()
	config.Load(*config_path)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := driver.NewDB()
	if err != nil {
		log.Fatalf("Can't connect to Database: %s\n", err)
	}
	storage := storage.NewStorage(db)
	service := service.NewService(storage)
	handler := handler.NewHandler(service)
	router := handler.GetRoute()

	server := &http.Server{
		Addr:    config.Koanf.MustString("address"),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
