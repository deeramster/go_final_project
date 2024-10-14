package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deeramster/go_final_project/auth"
	"github.com/deeramster/go_final_project/db"
	"github.com/deeramster/go_final_project/handlers"

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
	http.HandleFunc("/api/signin", handlers.HandleSignIn)
	http.HandleFunc("/api/task", auth.Middleware(handlers.HandleTask))
	http.HandleFunc("/api/tasks", auth.Middleware(handlers.HandleTasks))
	http.HandleFunc("/api/task/done", auth.Middleware(handlers.HandleTaskDone))
	http.HandleFunc("/api/nextdate", handlers.HandleNextDate)

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
