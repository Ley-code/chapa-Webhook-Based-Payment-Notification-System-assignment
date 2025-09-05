package main

import (
	"log"
	"net/http"

	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/handler"
	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/repository"
	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/usecase"
)

func main() {
	port:= "8080"
	// 1. Create dependencies starting from the inner layer.
	paymentRepository := repository.NewPaymentRepostiory()
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepository)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)

	// 2. Register the handler and start the server.
	http.Handle("/api/v1/payment", paymentHandler)

	log.Println("Payment Processor server starting on http://localhost:8080")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("FATAL: Could not start server: %v", err)
	}
}