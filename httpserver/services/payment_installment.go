package services

import (
	"context"
	"errors"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/responses"
	"net/http"
)

type paymentInstallmentSvc struct {
	repo repositories.PaymentInstallmentRepo
}

func NewPaymentInstallmentSvc(repo repositories.PaymentInstallmentRepo) PaymentInstallmentSvc {
	return &paymentInstallmentSvc{repo: repo}
}

func (p *paymentInstallmentSvc) GetInstallmentByCustomer(ctx context.Context, customerId uint) *responses.Response {

	paymentInstalment, err := p.repo.FindInstallmentByCustomerId(ctx, customerId)
	if err != nil || len(paymentInstalment) == 0 {
		return responses.ErrorResponse(responses.M_NOT_FOUND, http.StatusNotFound, errors.New("customer accepted loan request not found"))
	}
	var response []responses.PaymentInstallmentResponse
	for _, val := range paymentInstalment {
		response = append(response, responses.PaymentInstallmentResponse{
			Id:        val.Id,
			Amount:    val.Amount,
			DueDate:   val.DueDate,
			Status:    val.Status,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
		})
	}
	return responses.SuccessResponseWithData(responses.M_OK, http.StatusOK, response)
}
