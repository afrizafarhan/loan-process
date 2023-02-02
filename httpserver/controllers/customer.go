package controllers

import (
	"github.com/gin-gonic/gin"
	"loan_process/exceptions"
	"loan_process/httpserver/services"
	"strconv"
)

type customerController struct {
	svc services.CustomerSvc
}

func NewCustomerController(svc services.CustomerSvc) CustomerController {
	return &customerController{svc: svc}
}

func (c *customerController) GetCustomers(ctx *gin.Context) {
	response := c.svc.GetCustomers(ctx)
	WriteJsonResponse(ctx, response)
}

func (c *customerController) GetDetailCustomer(ctx *gin.Context) {
	customerId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		exceptions.ValidationError(ctx, err)
		return
	}
	response := c.svc.GetDetailCustomer(ctx, uint(customerId))
	WriteJsonResponse(ctx, response)
}
