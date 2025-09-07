package repository

import (
	"sync"

	"github.com/Ley-code/chapa-Webhook-Based-Payment-Notification-System-assignment/server/domain"
)

type PaymentRepository struct {
	payments map[string]*domain.Payment
	mutex    sync.RWMutex
}

func NewPaymentRepostiory() *PaymentRepository {
	return &PaymentRepository{
		payments: make(map[string]*domain.Payment),
	}
}

// create stores a new payment record.
func (pr *PaymentRepository) Create(payment *domain.Payment) error {
	pr.mutex.Lock() // this is actually pretty cool, it prevents other concurrent calls to simultaneously read and write. 
	defer pr.mutex.Unlock()
	pr.payments[payment.ID] = payment
	return nil
}

// update modifies an existing payment record's status.
func (pr *PaymentRepository) UpdateStatus(id, status string) error {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	if payment, ok := pr.payments[id]; ok {
		payment.Status = status
	}
	return nil
}