package controllers

import (
	"github.com/gin-gonic/gin"
)

type LoanApplicationController interface {
	Create(ctx *gin.Context)
	GetLoanApplications(ctx *gin.Context)
	ReapplyLoanApplication(ctx *gin.Context)
}
