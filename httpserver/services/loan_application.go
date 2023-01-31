package services

import (
	"context"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/request"
	"loan_process/httpserver/responses"
)

type loanApplicationSvc struct {
	customer repositories.CustomerRepo
}

func NewLoanApplicationSvc(customer repositories.CustomerRepo) *loanApplicationSvc {
	return &loanApplicationSvc{
		customer: customer,
	}
}

func (l loanApplicationSvc) CreateLoanApplication(ctx context.Context, application *request.CreateLoanApplication) *responses.ResponseMessage {
	//TODO implement me
	panic("implement me")
}
