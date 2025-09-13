package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io" 
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- INCOMING WEBHOOK RECEIVED (SIGNATURE VERIFIED) ---")
	body, _ := io.ReadAll(r.Body)
	log.Printf("Payload: %s", string(body))
	log.Println("-------------------------------------------------")
	w.WriteHeader(http.StatusOK)
}

func signatureVerificationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR: Cannot read body: %v", err)
			http.Error(w, "Cannot read body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		receivedSignature := r.Header.Get("X-Chapa-Signature")
		if receivedSignature == "" {
			log.Println("WARN: Missing X-Chapa-Signature header")
			http.Error(w, "Missing signature", http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("WEBHOOK_SECRET_KEY")
		if secret == "" {
			log.Printf("ERROR: WEBHOOK_SECRET_KEY is not set on client.")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h := hmac.New(sha256.New, []byte(secret))
		h.Write(body) // Hash the raw body bytes
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		if !hmac.Equal([]byte(receivedSignature), []byte(expectedSignature)) {
			log.Println("WARN: Invalid signature received")
			http.Error(w, "Invalid signature", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found")
		}
	}


    http.HandleFunc("/webhook", signatureVerificationMiddleware(webhookHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Mock Client listening for webhooks on http://localhost:%s/webhook", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("FATAL: Could not start client server: %v", err)
	}
}