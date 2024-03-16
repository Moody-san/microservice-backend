package main

import (
	"context"
	"errors"
	"github.com/AjaxAueleke/e-commerce/productService/api"
	"github.com/AjaxAueleke/e-commerce/productService/internal/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	rabbitmqConnectionString := os.Getenv("RABBITMQ_URL")
	rabbitMQSvc := service.NewRabbitMQService(rabbitmqConnectionString)
	productService := service.NewProductService(db)

	r := mux.NewRouter()
	api.RegisterProductRoutes(r, productService)

	http.Handle("/", r)
	log.Println("Product service started on :9020")
	go rabbitMQSvc.StartListeningForUserDeleteEvents(productService)
	httpServer := &http.Server{
		Addr:    ":9020",
		Handler: r,
	}

	go func() {
		log.Println("Product service started on :9020")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown failed: %v", err)
	}

	log.Println("Product service shutdown gracefully")
}
