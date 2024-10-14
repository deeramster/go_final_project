package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deeramster/go_final_project/auth"
	"github.com/deeramster/go_final_project/config"
	"github.com/deeramster/go_final_project/db"
	"github.com/deeramster/go_final_project/handlers"
)

func main() {
	//Get env variables
	config.LoadConfig()

	port := config.AppConfig.Port

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
