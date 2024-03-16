package main

import (
	"github.com/AjaxAueleke/e-commerce/paymentService/internal/model"
	"log"
	"net/http"

	"github.com/AjaxAueleke/e-commerce/paymentService/api"
	"github.com/AjaxAueleke/e-commerce/paymentService/internal/service"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Database connection string; adjust as per your MySQL setup
	dsn := "user:password@tcp(localhost:3306)/paymentservicedb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// AutoMigrate the `Payment` model; ensure it's defined and imported
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
