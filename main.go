package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//Get .env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default environment variables")
	}

	port := getPort()

	//Static files
	http.Handle("/", http.FileServer(http.Dir("./web")))

	//Init DB or create if not exist
	initDB()

	//Server initialization
	fmt.Printf("Server running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
