package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Yusuf-Doni/web-go-CRUD/database"
	"github.com/Yusuf-Doni/web-go-CRUD/routes"
	"github.com/Yusuf-Doni/web-go-CRUD/service"
	"github.com/joho/godotenv"
)

func main() {
	//Read file .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Tidak menemukan file .env, gunakan environment sistem")
	}

	// Initialize PostgreSQL database connection
	log.Println("Initializing PostgreSQL database connection...")
	db := database.InitDatabase()
	defer db.Close()

	// Create book service
	bookService := service.NewBookService(db)

	// Initialize database schema
	err := bookService.InitializeSheet()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
	}

	server := http.NewServeMux()

	routes.MapRoutes(server, bookService, db)

	log.Println("🚀 Server starting on :9000")
	log.Println("🔐 Login: http://localhost:9000/login")
	log.Println("📝 Register: http://localhost:9000/register")
	log.Println("📊 Dashboard: http://localhost:9000/dashboard")
	log.Println("➕ Add Product: http://localhost:9000/addproduct")
	log.Println("👤 Default Admin: username=admin, password=admin123")

	port := os.Getenv("SERVER_PORT")
	http.ListenAndServe(":"+port, server)
}
