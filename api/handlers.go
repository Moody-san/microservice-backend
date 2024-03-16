package api

import (
	"encoding/json"
	"errors"
	"github.com/AjaxAueleke/e-commerce/paymentService/internal/model"
	"github.com/AjaxAueleke/e-commerce/paymentService/internal/service"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func RegisterPaymentRoutes(r *mux.Router, s *service.PaymentService) {
	r.HandleFunc("/payments", createPaymentHandler(s)).Methods("POST")
	r.HandleFunc("/payments/{id}", getPaymentHandler(s)).Methods("GET")
	r.HandleFunc("/payments/{id}/status", updatePaymentStatusHandler(s)).Methods("PUT")
	r.HandleFunc("/payments", listPaymentsHandler(s)).Methods("GET")
}

// getPaymentHandler retrieves a payment by its ID
func getPaymentHandler(s *service.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid payment ID", http.StatusBadRequest)
			return
		}

		payment, err := s.GetPayment(uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Payment not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to fetch payment", http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(payment)
	}
}

func updatePaymentStatusHandler(s *service.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid payment ID", http.StatusBadRequest)
			return
		}

		var update struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := s.UpdatePaymentStatus(uint(id), update.Status); err != nil {
			http.Error(w, "Failed to update payment status", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Payment status updated successfully"})
	}
}

// createPaymentHandler creates a new payment record
func createPaymentHandler(s *service.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payment model.Payment
		if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := s.CreatePayment(&payment); err != nil {
			http.Error(w, "Failed to create payment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(payment)
	}
}

// listPaymentsHandler lists all payments
func listPaymentsHandler(s *service.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payments, err := s.ListPayments()
		if err != nil {
			http.Error(w, "Failed to fetch payments", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(payments)
	}
}
