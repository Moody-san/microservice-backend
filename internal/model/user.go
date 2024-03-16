package model

import "time"

// User represents a user in the system.
// @Description User object contains user information.
type User struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"` // Consider storing hashed passwords
	Email       string    `json:"email"`
	CompanyName string    `json:"companyName"`
	PhoneNumber string    `json:"phoneNumber"`
	Role        string    `json:"role"`
	Image       string    `json:"image"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}
