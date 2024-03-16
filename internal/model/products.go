package model

type Product struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Title    string  `json:"title"`
	Detail   string  `json:"detail"`
	Price    float64 `json:"price"`
	Image    string  `json:"image"` // Consider storing image URLs or identifiers for an external storage service instead
	UserID   uint    `json:"userID"`
	Quantity uint    `json:"quantity"`
}
