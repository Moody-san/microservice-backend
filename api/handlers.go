package api

import (
	"encoding/json"
	"github.com/AjaxAueleke/e-commerce/orderingService/internal/model"
	"github.com/AjaxAueleke/e-commerce/orderingService/internal/service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func RegisterOrderRoutes(r *mux.Router, s *service.OrderService) {
	r.HandleFunc("/orders", createOrderHandler(s)).Methods("POST")
	r.HandleFunc("/orders/{id}", getOrderHandler(s)).Methods("GET")
	r.HandleFunc("/orders/{id}", updateOrderHandler(s)).Methods("PUT")
	r.HandleFunc("/orders/{id}", deleteOrderHandler(s)).Methods("DELETE")
	r.HandleFunc("/orders", listOrdersHandler(s)).Methods("GET")
}

func createOrderHandler(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order model.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := s.CreateOrder(&order); err != nil {
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}

func getOrderHandler(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		order, err := s.GetOrder(uint(id))
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(order)
	}
}

func updateOrderHandler(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		var order model.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		order.ID = uint(id) // Ensure the order ID is set to the one from the URL.
		if err := s.UpdateOrder(&order); err != nil {
			http.Error(w, "Failed to update order", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(order)
	}
}

func deleteOrderHandler(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		if err := s.DeleteOrder(uint(id)); err != nil {
			http.Error(w, "Failed to delete order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // No content to send back on deletion
	}
}

func listOrdersHandler(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := s.ListOrders()
		if err != nil {
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(orders)
	}
}
