package controllers

import (
	"github.com/gin-gonic/gin"
)

type LoanApplicationController interface {
	Create(ctx *gin.Context)
}
