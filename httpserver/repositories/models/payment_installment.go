package models

import "time"

type PaymentInstallment struct {
	Id                    uint
	CustomerId            uint
	CustomerLoanRequestId uint
	Amount                uint
	DueDate               time.Time
	Status                string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
