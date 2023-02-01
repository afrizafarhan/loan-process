package models

import "time"

type CustomerLoanRequest struct {
	Id         uint
	CustomerId uint
	Amount     uint
	Tenor      uint
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
