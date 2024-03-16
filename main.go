package main

import (
	"fmt"
	"github.com/AjaxAueleke/e-commerce/paymentService/internal/model"
	"log"
	"net/http"
	"os"

	"github.com/AjaxAueleke/e-commerce/paymentService/api"
	"github.com/AjaxAueleke/e-commerce/paymentService/internal/service"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	log.Printf("dbUser: %v", dbUser)
	log.Printf("dbName: %v", dbName)
	log.Printf("dbPassword: %v", dbPassword)
	log.Printf("dbHost: %v", dbHost)
	log.Printf("dbPort: %v", dbPort)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	rabbitmqConnectionString := os.Getenv("RABBITMQ_URL")
	log.Printf("RABBIT_MQURL: %v", rabbitmqConnectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&model.Payment{})

	paymentService := service.NewPaymentService(db)
	router := mux.NewRouter()

	// Register payment routes
	api.RegisterPaymentRoutes(router, paymentService)

	// Starting the server
	log.Println("Starting PaymentService on port :9000...")
	if err := http.ListenAndServe(":9000", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
