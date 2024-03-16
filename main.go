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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can't read environment variables")
	}
	dsn := "newuser:user_password@tcp(localhost:3306)/productservicedb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	// Graceful shutdown handling
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
