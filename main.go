package main

import (
	"fmt"
	"github.com/AjaxAueleke/e-commerce/userService/api"
	"github.com/AjaxAueleke/e-commerce/userService/internal/service"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	dbUser := "newuser"
	dbPassword := "user_password"
	dbName := "userservicedb"

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(host.docker.internal:3306)/%s?parseTime=true", dbUser, dbPassword, dbName)), &gorm.Config{})
	if err != nil {
		log.Fatal("Can't connect to the database")
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
		log.Println("Error creating the server")
		return
	}
}
