package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
		return
	}

	fmt.Println("INCOMING WEBHOOK RECEIVED")
	fmt.Printf("Payment ID: %s", payload["id"])
	fmt.Printf("Status: %s", payload["status"])
	fmt.Printf("Amount: %.2f %s", payload["amount"], payload["currency"])
	fmt.Println("-------------------------------")

	w.WriteHeader(http.StatusOK)
}

func main() {
	port:= "8081"
	http.HandleFunc("/webhook", webhookHandler)
	fmt.Println("Mock Client listening for webhooks on http://localhost:8081/webhook")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("FATAL: Could not start client server: %v\n", err)
	}
}