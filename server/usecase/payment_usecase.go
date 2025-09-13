package usecase

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/server/domain"
	"github.com/google/uuid"
)

// PaymentRepositoryInterface defines the contract for our data layer.
// The usecase depends on this interface, not the concrete implementation.
type PaymentRepositoryInterface interface {
	Create(payment *domain.Payment) error
	UpdateStatus(id, status string) error
}

type PaymentUsecase struct {
	repository PaymentRepositoryInterface
}

func NewPaymentUsecase(repo PaymentRepositoryInterface) *PaymentUsecase {
	return &PaymentUsecase{
		repository: repo,
	}
}

func (pu *PaymentUsecase) ProcessPayment(paymentReq domain.PaymentRequest) error {
	// create the internal payment object
	newPayment := &domain.Payment{
		ID:         uuid.New().String(),
		Status:     "PENDING",
		Amount:     paymentReq.Amount,
		Currency:   paymentReq.Currency,
		WebhookURL: paymentReq.WebhookURL,
	}

	// create the initial "PENDING" state in memory using hashmap
	if err := pu.repository.Create(newPayment); err != nil {
		log.Printf("ERROR: Failed to create payment: %v", err)
		return err
	}

	// process asynchronously so we can return a response to the user immediately.
	// concurrency baby:)
	go pu.simulateAndNotify(newPayment)
	return nil
}
func (pu *PaymentUsecase) simulateAndNotify(payment *domain.Payment) {
	log.Printf("INFO: Processing payment %s...", payment.ID)
	
	//simulating processing payment
	time.Sleep(3 * time.Second)

	// marking payment as processed and updating it in the repository.
	newStatus := "PROCESSED"
	log.Printf("INFO: Payment %s is now %s.", payment.ID, newStatus)
	pu.repository.UpdateStatus(payment.ID, newStatus)
	payment.Status = newStatus 

	payload, err := json.Marshal(payment)
	if err != nil {
		log.Printf("ERROR: Failed to create webhook payload for payment %s: %v", payment.ID, err)
		return
	}

	secret := os.Getenv("WEBHOOK_SECRET_KEY")
	if secret == "" {
		log.Printf("ERROR: WEBHOOK_SECRET_KEY is not set. Cannot send webhook.")
		return
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	signature := hex.EncodeToString(h.Sum(nil))

	req, err := http.NewRequest("POST", payment.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("ERROR: Failed to create webhook request for payment %s: %v", payment.ID, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Chapa-Signature", signature) 

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to send webhook for payment %s: %v", payment.ID, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("INFO: Webhook for payment %s sent successfully. Status: %s", payment.ID, resp.Status)
}