package models

import "time"

type DailyLoanRequest struct {
	Id                 uint
	CurrentDateRequest time.Time
	Request            uint
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
