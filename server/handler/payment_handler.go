package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/domain"
	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/usecase"
)

type PaymentHandler struct {
	usecase *usecase.PaymentUsecase
}

func NewPaymentHandler(uc *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		usecase: uc,
	}
}

// ServeHTTP makes our PaymentHandler compatible with the standard http.Handler interface. just found out today actually
func (ph *PaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var paymentRequest domain.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// pass the request to the business logic layer
	err := ph.usecase.ProcessPayment(paymentRequest)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// respond immediately with "Accepted"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Payment processing initiated"})
}