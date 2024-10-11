package main

import (
	"fmt"
	"github.com/deeramster/go_final_project/db"
	"github.com/deeramster/go_final_project/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//Get env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading env file, using default environment variables")
	}

	port := getPort()

	//Static files
	http.Handle("/", http.FileServer(http.Dir("./web")))

	//API endpoints
	http.HandleFunc("/api/task", handlers.HandleTask)
	http.HandleFunc("/api/tasks", handlers.HandleTasks)
	http.HandleFunc("/api/task/done", handlers.HandleTaskDone)

	//Init DB or create if not exist
	db.InitDB()

	//Server initialization
	fmt.Printf("Server running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
