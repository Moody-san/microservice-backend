package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AjaxAueleke/e-commerce/orderingService/api"
	"github.com/AjaxAueleke/e-commerce/orderingService/internal/service"
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
	orderService := service.NewOrderService(db)
	r := mux.NewRouter()
	api.RegisterOrderRoutes(r, orderService)
	http.Handle("/", r)
	httpServer := &http.Server{
		Addr:    ":9020",
		Handler: r,
	}
	go func() {
		log.Println("Order service started on :9020")
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
	log.Println("Order service shutdown gracefully")

}
