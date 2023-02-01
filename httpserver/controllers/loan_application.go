package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"loan_process/httpserver/request"
	"loan_process/httpserver/services"
	"net/http"
)

type loanApplicationController struct {
	svc services.LoanApplicationSvc
}

func NewLoanApplicationController(svc services.LoanApplicationSvc) LoanApplicationController {
	return &loanApplicationController{
		svc: svc,
	}
}

func (l *loanApplicationController) Create(ctx *gin.Context) {
	var req request.CreateLoanApplication
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = validator.New().Struct(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}
	response := l.svc.CreateLoanApplication(ctx, &req)
	WriteJsonResponse(ctx, response)
}
