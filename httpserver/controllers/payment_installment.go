package controllers

import (
	"github.com/gin-gonic/gin"
	"loan_process/exceptions"
	"loan_process/httpserver/services"
	"strconv"
)

type paymentInstallmentController struct {
	svc services.PaymentInstallmentSvc
}

func NewPaymentInstallmentController(svc services.PaymentInstallmentSvc) PaymentInstallmentController {
	return &paymentInstallmentController{svc: svc}
}

func (p *paymentInstallmentController) GetInstallmentByCustomerAndLoanRequest(ctx *gin.Context) {
	customerId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		exceptions.ValidationError(ctx, err)
		return
	}
	response := p.svc.GetInstallmentByCustomer(ctx, uint(customerId))
	WriteJsonResponse(ctx, response)
}
