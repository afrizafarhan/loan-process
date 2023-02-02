package services

import (
	"context"
	"loan_process/httpserver/repositories"
	"time"
)

type dailyLoanRequestSvc struct {
	dailyLoan repositories.DailyLoanRequestRepo
}

func NewDailyLoanRequestSvc(dailyLoan repositories.DailyLoanRequestRepo) DailyLoanRequestSvc {
	return &dailyLoanRequestSvc{dailyLoan: dailyLoan}
}

func (d *dailyLoanRequestSvc) CheckDailyLoanRequest(ctx context.Context) int {
	dailyLoanRequest, err := d.dailyLoan.FindDailyLoanRequestByDate(ctx, time.Now())
	if err != nil {
		return 0
	}
	return int(dailyLoanRequest.Request)
}
