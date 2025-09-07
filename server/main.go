package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/server/handler"
	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/server/repository"
	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/server/usecase"
	"github.com/joho/godotenv" 
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found")
		}
	}

	paymentRepository := repository.NewPaymentRepostiory()
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepository)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)

	http.Handle("/api/v1/payment", paymentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Printf("Payment Processor server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("FATAL: Could not start server: %v", err)
	}
}