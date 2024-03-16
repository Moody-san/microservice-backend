package model

import "time"

type Payment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // e.g., "pending", "completed", "failed"
	CreatedAt time.Time `json:"created_at"`
}
