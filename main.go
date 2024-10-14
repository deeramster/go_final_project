package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/deeramster/go_final_project/auth"
	"github.com/deeramster/go_final_project/config"
	"github.com/deeramster/go_final_project/db"
	"github.com/deeramster/go_final_project/handlers"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Get the configured port from the config
	port := config.AppConfig.Port

	// Initialize the database connection
	db.InitDB()

	// Ensure the database is closed when the application exits
	defer db.CloseDB()

	// Static files (serves files from the "web" directory)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// API Endpoints
	http.HandleFunc("/api/signin", handlers.HandleSignIn)
	http.HandleFunc("/api/task", auth.Middleware(handlers.HandleTask))
	http.HandleFunc("/api/tasks", auth.Middleware(handlers.HandleTasks))
	http.HandleFunc("/api/task/done", auth.Middleware(handlers.HandleTaskDone))
	http.HandleFunc("/api/nextdate", handlers.HandleNextDate)

	// Start the server in a separate goroutine
	go func() {
		fmt.Printf("Server running on port %s...\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal("Error starting server:", err)
		}
	}()

	// Graceful shutdown on receiving a termination signal (Ctrl+C, system signal)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan

	// Graceful shutdown message
	log.Println("Shutting down server...")
}
