package services

import (
	"context"
	"loan_process/httpserver/repositories/models"
	"loan_process/httpserver/request"
	"loan_process/httpserver/responses"
)

type LoanApplicationSvc interface {
	CreateLoanApplication(ctx context.Context, application *request.CreateLoanApplication) *responses.Response
	GetLoanApplication(ctx context.Context) *responses.Response
	ReapplyLoanApplication(ctx context.Context, customerId uint, application *request.ReapplyLoanApplication) *responses.Response
	createCustomerLoanRequest(ctx context.Context, customer *models.Customer, loanAmount uint, tenor uint) error
}

type DailyLoanRequestSvc interface {
	CheckDailyLoanRequest(ctx context.Context) int
}
