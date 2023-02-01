package services

import (
	"context"
	"loan_process/httpserver/request"
	"loan_process/httpserver/responses"
)

type LoanApplicationSvc interface {
	CreateLoanApplication(ctx context.Context, application *request.CreateLoanApplication) *responses.Response
}
