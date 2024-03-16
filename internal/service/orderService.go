package service

import (
	"github.com/AjaxAueleke/e-commerce/orderingService/internal/model"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) CreateOrder(order *model.Order) error {
	return s.db.Create(order).Error
}

func (s *OrderService) GetOrder(id uint) (*model.Order, error) {
	var order model.Order
	err := s.db.First(&order, id).Error
	return &order, err
}

func (s *OrderService) UpdateOrder(order *model.Order) error {
	return s.db.Save(order).Error
}

func (s *OrderService) DeleteOrder(id uint) error {
	return s.db.Delete(&model.Order{}, id).Error
}

func (s *OrderService) ListOrders() ([]model.Order, error) {
	var orders []model.Order
	err := s.db.Find(&orders).Error
	return orders, err
}
