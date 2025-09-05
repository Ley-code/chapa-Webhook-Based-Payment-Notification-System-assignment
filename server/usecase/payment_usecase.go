package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/domain"
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
		fmt.Printf("ERROR: Failed to create payment: %v", err)
		return err
	}

	// process asynchronously so we can return a response to the user immediately.
	// concurrency baby:)
	go pu.simulateAndNotify(newPayment)
	return nil
}
func (pu *PaymentUsecase) simulateAndNotify(payment *domain.Payment) {
	fmt.Printf("INFO: Processing payment %s...", payment.ID)
	
	//simulating processing payment
	time.Sleep(3 * time.Second)

	// marking payment as processed and updating it in the repository.
	newStatus := "PROCESSED"
	fmt.Printf("INFO: Payment %s is now %s.", payment.ID, newStatus)
	pu.repository.UpdateStatus(payment.ID, newStatus)
	payment.Status = newStatus 

	payload, err := json.Marshal(payment)
	if err != nil {
		fmt.Printf("ERROR: Failed to create webhook payload for payment %s: %v", payment.ID, err)
		return
	}

	// this is where the magic happens. sending the webhook notification.
	resp, err := http.Post(payment.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("ERROR: Failed to send webhook for payment %s: %v", payment.ID, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("INFO: Webhook for payment %s sent successfully. Status: %s", payment.ID, resp.Status)
}