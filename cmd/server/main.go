package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"lwnra-devo-api/config"
	"lwnra-devo-api/database"
	"lwnra-devo-api/facebook"
	"lwnra-devo-api/handlers"
	"lwnra-devo-api/middleware"
	"lwnra-devo-api/routes"
	"lwnra-devo-api/scheduler"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.FacebookToken == "" {
		log.Println("Warning: FB_ACCESS_TOKEN not set. Facebook sync will not work.")
	}

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Facebook client
	fbClient := facebook.New(cfg.FacebookToken)

	// Initialize scheduler
	sched := scheduler.New(db, fbClient)
	
	// Start scheduler if Facebook token is available
	if cfg.FacebookToken != "" {
		sched.Start()
		defer sched.Stop()
		
		nextRun := sched.GetNextRun()
		if !nextRun.IsZero() {
			fmt.Printf("⏰ Next sync scheduled for: %s\n", nextRun.Format("2006-01-02 15:04:05 MST"))
		}
	} else {
		fmt.Println("⚠️  Scheduler disabled - FB_ACCESS_TOKEN not set")
	}

	// Initialize handlers
	devotionalHandler := handlers.NewDevotionalHandler(db, fbClient)
	systemHandler := handlers.NewSystemHandler(sched)

	// Initialize router
	router := routes.NewRouter(devotionalHandler, systemHandler)

	// Apply middleware
	handler := middleware.Logger(middleware.Recovery(middleware.CORS(router)))

	// Start server
	fmt.Printf("🚀 LWNRA Devotional API v1.0.0\n")
	fmt.Printf("📖 Environment: %s\n", cfg.Environment)
	fmt.Printf("🌐 Server starting on port %s\n", cfg.Port)
	fmt.Printf("📄 API Docs: http://localhost:%s/\n", cfg.Port)
	fmt.Printf("🏥 Health: http://localhost:%s/health\n", cfg.Port)
	fmt.Printf("📖 Devotionals: http://localhost:%s/api/devotionals\n", cfg.Port)
	fmt.Printf("⏰ Scheduler: http://localhost:%s/api/scheduler/status\n", cfg.Port)
	fmt.Println("Press Ctrl+C to stop the server")

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\n🛑 Shutting down gracefully...")
	
	// Stop scheduler
	if cfg.FacebookToken != "" {
		sched.Stop()
	}
	
	fmt.Println("✅ Server stopped")
}
