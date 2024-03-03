package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AjaxAueleke/e-commerce/userService/internal/model"
	"github.com/AjaxAueleke/e-commerce/userService/internal/service"
)

func CreateUserHandler(s *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Error decoding user data: %v", err)
			http.Error(w, "Invalid user data", http.StatusBadRequest)
			return
		}

		if err := s.CreateUser(&user); err != nil {
			log.Printf("Error creating user: %v", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
	}
}

func UpdateUserHandler(s *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, ok := vars["id"]
		if !ok {
			log.Println("User ID is missing in URL path")
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Error decoding user data: %v", err)
			http.Error(w, "Invalid user data", http.StatusBadRequest)
			return
		}
		userID_int, err_userId := strconv.Atoi(userID)
		if err_userId != nil {
			log.Printf("Error converting userid: %v", err_userId)
			http.Error(w, "Invalid user id", http.StatusBadRequest)
		}
		user.ID = uint(userID_int) // Assuming ID is a string. Convert or adapt as necessary

		if err := s.UpdateUser(&user); err != nil {
			log.Printf("Error updating user: %v", err)
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
	}
}

func DeleteUserHandler(s *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, ok := vars["id"]
		if !ok {
			log.Println("User ID is missing in URL path")
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}
		userID_int, err := strconv.Atoi(userID)
		if !ok {
			log.Println("User ID is not a valid integer")
			http.Error(w, "User ID should be an integer", http.StatusBadRequest)
			return
		}

		if err := s.DeleteUser(uint(userID_int)); err != nil {
			log.Printf("Error deleting user: %v", err)
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
		if err != nil {
			return
		}
	}
}
func ListUsersHandler(s *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.ListUsers()
		if err != nil {
			log.Printf("Error fetching users: %v", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("Error encoding users response: %v", err)
			http.Error(w, "Error processing users data", http.StatusInternalServerError)
		}
	}
}

func LoginHandler(s *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := s.Authenticate(credentials.Email, credentials.Password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := s.GenerateJWT(user.ID)

		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

var jwtKey = []byte("a high stake security key")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
