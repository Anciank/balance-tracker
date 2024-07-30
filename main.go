package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"balance-tracker/handlers"
	"balance-tracker/repositories"
	"balance-tracker/services"
	"balance-tracker/utils"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Load environment variables
	envs := utils.NewEnvEngine()
	dbHost := envs.LoadEnv("DB_HOST")
	dbPort := envs.LoadEnv("DB_PORT")
	dbUser := envs.LoadEnv("DB_USER")
	dbPassword := envs.LoadEnv("DB_PASSWORD")
	dbName := envs.LoadEnv("DB_NAME")

	// Connect to PostgreSQL database
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create repositories
	balanceRepository := repositories.NewBalanceRepository(db)
	userRepository := repositories.NewUserRepository(db)
	sessionRepository := repositories.NewSessionRepository(db)

	// Create services
	authService := services.NewAuthService(userRepository, sessionRepository)
	balanceService := services.NewBalanceService(balanceRepository)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	balanceHandler := handlers.NewBalanceHandler(balanceService)
	pageHandler := handlers.NewPageHandler(balanceService)

	// Create HTTP server
	server := http.NewServeMux()

	// Register handlers
	server.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			pageHandler.HandleLoginPage(w, r)
		} else if r.Method == "POST" {
			authHandler.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			pageHandler.HandleRegisterPage(w, r)
		} else if r.Method == "POST" {
			authHandler.Register(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server.HandleFunc("/logout", authHandler.Logout)

	// Authenticated routes
	authMiddleware := authHandler.AuthMiddleware()

	server.HandleFunc("/balances", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			balanceHandler.GetBalances(w, r)
		case http.MethodPost:
			balanceHandler.CreateBalance(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	server.HandleFunc("/balances/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			balanceHandler.GetBalance(w, r)
		case http.MethodPost:
			balanceHandler.CreateBalance(w, r)
		case http.MethodPut:
			balanceHandler.UpdateBalance(w, r)
		case http.MethodDelete:
			balanceHandler.DeleteBalance(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	server.HandleFunc("/", authMiddleware(pageHandler.HandleIndexPage))
	server.HandleFunc("/htmx.min.js", pageHandler.HandleHtmxServe)
	server.HandleFunc("/tailwind.js", pageHandler.HandleTailwindServe)
	server.HandleFunc("/static/*", pageHandler.HandleStaticServe)
	server.HandleFunc("/create-transaction", authMiddleware(balanceHandler.CreateTransaction))

	// Start HTTP server
	log.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", server)
}
