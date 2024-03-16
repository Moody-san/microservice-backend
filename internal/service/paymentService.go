package service

import (
	"github.com/AjaxAueleke/e-commerce/paymentService/internal/model"
	"gorm.io/gorm"
)

type PaymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{db: db}
}

func (s *PaymentService) CreatePayment(payment *model.Payment) error {
	return s.db.Create(payment).Error
}

func (s *PaymentService) GetPayment(id uint) (*model.Payment, error) {
	var payment model.Payment
	err := s.db.First(&payment, id).Error
	return &payment, err
}

func (s *PaymentService) UpdatePaymentStatus(id uint, status string) error {
	return s.db.Model(&model.Payment{}).Where("id = ?", id).Update("status", status).Error
}

func (s *PaymentService) ListPayments() ([]model.Payment, error) {
	var payments []model.Payment
	err := s.db.Find(&payments).Error
	return payments, err
}
