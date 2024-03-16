package service

import (
	"github.com/AjaxAueleke/e-commerce/productService/internal/model"
	"gorm.io/gorm"
	"strings"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	err := db.AutoMigrate(&model.Product{})
	if err != nil {
		panic("failed to auto migrate")
	}
	return &ProductService{db: db}
}

func (s *ProductService) CreateProduct(product *model.Product) error {
	return s.db.Create(product).Error
}

func (s *ProductService) GetProduct(id uint) (*model.Product, error) {
	var product model.Product
	err := s.db.First(&product, id).Error
	return &product, err
}

func (s *ProductService) UpdateProduct(product *model.Product) error {
	return s.db.Save(product).Error
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.db.Delete(&model.Product{}, id).Error
}

func (s *ProductService) ListProducts(page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var count int64
	db := s.db.Where("quantity > ?", 0) // Filter out products with quantity 0
	db.Model(&model.Product{}).Count(&count)

	// Calculate the offset based on the page and pageSize
	offset := (page - 1) * pageSize

	// First, find the total number of products to set up pagination on the client side
	result := db.Model(&model.Product{}).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Then, retrieve a paginated list of products
	err := db.Offset(offset).Limit(pageSize).Find(&products).Error
	return products, count, err
}

func (s *ProductService) DeleteProductsByUserID(userID uint) error {
	result := s.db.Where("user_id = ?", userID).Delete(&model.Product{})
	return result.Error
}

func (s *ProductService) SearchProducts(query string, sort string, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var count int64

	db := s.db.Model(&model.Product{}).Where("quantity > ?", 0)

	if query != "" {
		query = "%" + strings.ToLower(query) + "%"
		db = db.Where("LOWER(title) LIKE ?", query)
	}

	// Define a map of valid sort fields to prevent SQL injection
	validSortFields := map[string]bool{
		"title asc":  true,
		"title desc": true,
		"price asc":  true,
		"price desc": true,
	}

	if sort != "" {
		// Check if the sort parameter is valid
		if _, ok := validSortFields[sort]; ok {
			db = db.Order(sort)
		}
	}

	// Pagination
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)

	db = db.Find(&products).Count(&count)

	err := db.Error
	return products, count, err
}

func (s *ProductService) GetProductsByUserID(userID uint, page, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var count int64

	// First, count the total number of products for the given user ID
	db := s.db.Where("user_id = ?", userID)
	db.Model(&model.Product{}).Count(&count)

	// Apply pagination and retrieve the products
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&products).Error

	return products, count, err
}
