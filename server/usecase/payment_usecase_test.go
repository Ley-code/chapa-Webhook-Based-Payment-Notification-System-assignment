// in processor/usecase/payment_usecase_test.go

package usecase

import (
	"fmt"
	"testing"

	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/server/domain"
)

type MockPaymentRepository struct {
	CreateFunc       func(payment *domain.Payment) error
	UpdateStatusFunc func(id, status string) error

	CreateCalled bool
	LastPayment  *domain.Payment 
}

func (m *MockPaymentRepository) Create(payment *domain.Payment) error {
	m.CreateCalled = true
	m.LastPayment = payment
	if m.CreateFunc != nil {
		return m.CreateFunc(payment)
	}
	return nil
}

func (m *MockPaymentRepository) UpdateStatus(id, status string) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(id, status)
	}
	return nil
}

// TestProcessPayment verifies the core logic of creating a payment i.e spying
func TestProcessPayment(t *testing.T) {
	// arrange
	mockRepo := &MockPaymentRepository{}
	paymentUsecase := NewPaymentUsecase(mockRepo)

	request := domain.PaymentRequest{
		Amount:     100.00,
		Currency:   "USD",
		WebhookURL: "http://test.com/hook",
	}

	//act
	err := paymentUsecase.ProcessPayment(request)

	//assert

	if err != nil {
		t.Errorf("ProcessPayment() returned an unexpected error: %v", err)
	}

	if !mockRepo.CreateCalled {
		t.Error("expected repository.Create() to be called, but it wasn't")
	}

	if mockRepo.LastPayment == nil {
		t.Fatal("repository.Create() was called with a nil payment")
	}
	if mockRepo.LastPayment.Status != "PENDING" {
		t.Errorf("expected payment status to be 'PENDING', but got '%s'", mockRepo.LastPayment.Status)
	}

	if mockRepo.LastPayment.ID == "" {
		t.Error("expected payment ID to be generated, but it was empty")
	}

	if mockRepo.LastPayment.Amount != request.Amount {
		t.Errorf("expected amount to be %f, but got %f", request.Amount, mockRepo.LastPayment.Amount)
	}
}

func TestProcessPayment_FailsOnCreate(t *testing.T) {

	//arrange
	mockRepo := &MockPaymentRepository{}
	mockRepo.UpdateStatusFunc = func(id, status string) error {
		return fmt.Errorf("failed to update payment status")
	}
	paymentUsecase := NewPaymentUsecase(mockRepo)
	testPaymentRequest := domain.PaymentRequest{
		Amount:     100.00,
		Currency:   "USD",
		WebhookURL: "http://test.com/hook",
	}
	//act
	err := paymentUsecase.ProcessPayment(testPaymentRequest)
	
	//assert
	if err!=nil{
		t.Error("expected processpayment to return an error")
	}

}

func TestProcessPayment_FailsOnUpdate(t *testing.T) {

	//arrange
	mockRepo := &MockPaymentRepository{}
	mockRepo.UpdateStatusFunc = func(id, status string) error {
		return fmt.Errorf("failed to update payment status")
	}
	paymentUsecase := NewPaymentUsecase(mockRepo)
	testPaymentRequest := domain.PaymentRequest{
		Amount:     100.00,
		Currency:   "USD",
		WebhookURL: "http://test.com/hook",
	}
	//act
	err := paymentUsecase.ProcessPayment(testPaymentRequest)
	
	//assert
	if err!=nil{
		t.Error("expected processpayment to return an error")
	}

}

