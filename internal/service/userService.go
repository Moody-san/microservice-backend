package service

import (
	"errors"
	"fmt"
	"github.com/AjaxAueleke/e-commerce/userService/internal/model"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type UserService struct {
	db          *gorm.DB
	rabbitMQSvc *RabbitMQService
}

func NewUserService(db *gorm.DB, rabbitMQSvc *RabbitMQService) *UserService {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		panic("failed to auto migrate")
	}
	return &UserService{
		db:          db,
		rabbitMQSvc: rabbitMQSvc,
	}
}

func (s *UserService) CreateUser(user *model.User) error {
	var count int64
	s.db.Model(&model.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return errors.New("email already in use")
	}

	if user.Role == "Admin" {
		var adminCount int64
		s.db.Model(&model.User{}).Where("role = ?", "Admin").Count(&adminCount)
		if adminCount > 0 {
			return errors.New("an admin already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	result := s.db.Create(user)
	return result.Error
}

func (s *UserService) UpdateUser(user *model.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		var existingUser model.User
		if err := s.db.First(&existingUser, user.ID).Error; err != nil {
			return err
		}
		user.Password = existingUser.Password
	}

	result := s.db.Save(user)
	return result.Error
}

func (s *UserService) DeleteUser(userID uint) error {
	result := s.db.Delete(&model.User{}, userID)
	if result.Error != nil {
		return result.Error
	}
	eventType := "user.deleted"
	payload := []byte(fmt.Sprintf(`{"userID": %d, "message": "User deleted successfully."}`, userID))
	exchange := "user.deleted"
	routingKey := ""

	err := s.rabbitMQSvc.PublishEvent(exchange, routingKey, eventType, payload)
	if err != nil {
		log.Printf("Failed to publish user deletion event: %v", err)
	}
	return nil
}

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
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
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
