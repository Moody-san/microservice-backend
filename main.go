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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	log.Printf("dbUser: %v", dbUser)
	log.Printf("dbName: %v", dbName)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't connect to the database", err))
	}
	userService := service.NewUserService(db)

	r := mux.NewRouter()

	r.HandleFunc("/login", api.LoginHandler(userService)).Methods("POST")
	r.HandleFunc("/users", api.ListUsersHandler(userService)).Methods("GET")
	r.HandleFunc("/users", api.CreateUserHandler(userService)).Methods("POST")
	r.Handle("/users/{id}", api.JWTMiddleware(api.UpdateUserHandler(userService))).Methods("PUT")
	r.Handle("/users/{id}", api.JWTMiddleware(api.DeleteUserHandler(userService))).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	http.Handle("/", r)
	log.Println("Server started on :9000")
	err = http.ListenAndServe(":9000", nil)

	if err != nil {
		log.Fatal("Error creating the server")
		return
	}
}
