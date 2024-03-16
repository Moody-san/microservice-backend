package main

import (
	"fmt"
	"github.com/AjaxAueleke/e-commerce/userService/api"
	"github.com/AjaxAueleke/e-commerce/userService/internal/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the required details
		start := time.Now()
		log.Printf("Started %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)

		next.ServeHTTP(w, r)

		// You can also log the response status and the time taken to serve the request
		log.Printf("Completed in %v", time.Since(start))
	})
}
func main() {

	rabbitmqConnectionString := os.Getenv("RABBITMQ_URL")
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
		log.Fatal(fmt.Sprintf("Can't connect to the database", err))
	}

	rabbitMQSvc := service.NewRabbitMQService(rabbitmqConnectionString)
	err = rabbitMQSvc.SetupQueueAndBind("user.deleted", "userDeleteQueue")
	if err != nil {
		log.Fatalf("Failed to setup queue and bind it: %v", err)
	}
	userService := service.NewUserService(db, rabbitMQSvc)

	r := mux.NewRouter()

	r.HandleFunc("/login", api.LoginHandler(userService)).Methods("POST")
	r.HandleFunc("/users", api.ListUsersHandler(userService)).Methods("GET")
	r.HandleFunc("/users", api.CreateUserHandler(userService)).Methods("POST")
	r.Handle("/users/{id}", api.JWTMiddleware(api.UpdateUserHandler(userService))).Methods("PUT")
	r.Handle("/users/{id}", api.JWTMiddleware(api.DeleteUserHandler(userService))).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	http.Handle("/", LoggingMiddleware(r))

	log.Println("Server started on :9000")

	err = http.ListenAndServe(":9000", nil)

	if err != nil {
		log.Fatal("Error creating the server")
		return
	}
}
