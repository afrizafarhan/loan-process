package models

import "time"

type PaymentInstallment struct {
	Id         uint
	CustomerId uint
	Amount     uint
	DueDate    time.Time
	status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
