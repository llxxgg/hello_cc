package main

import (
	"fmt"
	"hello_cc/internal/user/handlers"
	"hello_cc/internal/user/services"
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialize services
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	// Define routes
	http.HandleFunc("/api/v1/users/register", userHandler.RegisterHandler)
	http.HandleFunc("/api/v1/users/login", userHandler.LoginHandler)
	http.HandleFunc("/api/v1/users/profile", userHandler.ProfileHandler)
	http.HandleFunc("/api/v1/users/usernames", userHandler.AllUsernamesHandler)

	// Default port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on :%s\n", port)
	fmt.Println("Endpoints:")
	fmt.Println("  POST   /api/v1/users/register - Register a new user")
	fmt.Println("  POST   /api/v1/users/login    - Login a user")
	fmt.Println("  POST   /api/v1/users/profile  - Get user profile")
	fmt.Println("  GET    /api/v1/users/usernames- Get all usernames")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
