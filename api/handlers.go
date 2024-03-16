package api

import (
	"encoding/json"
	"github.com/AjaxAueleke/e-commerce/productService/internal/model"
	"github.com/AjaxAueleke/e-commerce/productService/internal/service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func RegisterProductRoutes(r *mux.Router, s *service.ProductService) {
	r.HandleFunc("/products", createProductHandler(s)).Methods("POST")
	r.HandleFunc("/products/{id}", getProductHandler(s)).Methods("GET")
	r.HandleFunc("/products/{id}", updateProductHandler(s)).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProductHandler(s)).Methods("DELETE")
	r.HandleFunc("/products", listProductsHandler(s)).Methods("GET")
	r.HandleFunc("/users/{userID}/products", getProductsByUserIDHandler(s)).Methods("GET") // New route

}
func createProductHandler(s *service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product model.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := s.CreateProduct(&product); err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}
func getProductHandler(s *service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		product, err := s.GetProduct(uint(id))
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(product)
	}
}
func updateProductHandler(s *service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var product model.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		product.ID = uint(id) // Ensure the product ID is set to the one from the URL.
		if err := s.UpdateProduct(&product); err != nil {
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(product)
	}
}
func deleteProductHandler(s *service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		if err := s.DeleteProduct(uint(id)); err != nil {
			http.Error(w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content is appropriate for a DELETE success response
	}
}

func listProductsHandler(s *service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse optional query parameters
		query := r.URL.Query().Get("query")
		sort := r.URL.Query().Get("sort")
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("pageSize")

		// Set default values for pagination if not specified
		page, pageSize := 1, 10 // Default values
		var err error
		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil {
				http.Error(w, "Invalid page number", http.StatusBadRequest)
				return
			}
		}
		if pageSizeStr != "" {
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil {
				http.Error(w, "Invalid page size", http.StatusBadRequest)
				return
			}
		}

		var products []model.Product
		var total int64

		// If a query is provided, use the search method, otherwise list all products
		if query != "" {
			products, total, err = s.SearchProducts(query, sort, page, pageSize)
		} else {
			products, total, err = s.ListProducts(page, pageSize)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare and send the response
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"products": products,
			"total":    total,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func getProductsByUserIDHandler(s *service.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["userID"]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Parse optional query parameters for pagination
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("pageSize")
		page, pageSize := 1, 10 // Default values
		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil {
				http.Error(w, "Invalid page number", http.StatusBadRequest)
				return
			}
		}
		if pageSizeStr != "" {
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil {
				http.Error(w, "Invalid page size", http.StatusBadRequest)
				return
			}
		}

		products, total, err := s.GetProductsByUserID(uint(userID), page, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare and send the response
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"products": products,
			"total":    total,
		}
		json.NewEncoder(w).Encode(response)
	}
}
