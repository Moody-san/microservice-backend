package service

import (
	"errors"
	"fmt"
	"github.com/AjaxAueleke/e-commerce/userService/internal/model"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	// Automatically migrate your schema
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		panic("failed to auto migrate")
	}
	return &UserService{db: db}
}

// CreateUser inserts a new user into the database with hashed password and checks for email uniqueness
func (s *UserService) CreateUser(user *model.User) error {
	// Check for duplicate email
	var count int64
	s.db.Model(&model.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return errors.New("email already in use")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create user
	result := s.db.Create(user)
	return result.Error
}

func (s *UserService) UpdateUser(user *model.User) error {
	// If a new password is provided, hash it
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		// If no new password is provided, retain the existing one
		var existingUser model.User
		if err := s.db.First(&existingUser, user.ID).Error; err != nil {
			return err
		}
		user.Password = existingUser.Password
	}

	// Update user, GORM Save method updates all fields
	result := s.db.Save(user)
	return result.Error
}

// DeleteUser removes a user from the database
func (s *UserService) DeleteUser(userID uint) error {
	result := s.db.Delete(&model.User{}, userID)
	return result.Error
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user is not found
		}
		return nil, result.Error
	}
	return &user, nil
}

func (s *UserService) ListUsers() ([]model.User, error) {
	var users []model.User
	result := s.db.Model(&model.User{}).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	for i := range users {
		users[i].Password = "" // Clear passwords
	}
	return users, nil
}

func (s *UserService) Authenticate(email, password string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found
			return nil, errors.New("incorrect email or password")
		}
		// Database error
		return nil, err
	}

	// Check if the provided password matches the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Password does not match
		return nil, errors.New("incorrect email or password")
	}

	return &user, nil
}

var jwtKey = []byte("a high stake security key")

func (s *UserService) GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}
