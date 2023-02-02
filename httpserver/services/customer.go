package services

import (
	"context"
	"errors"
	"fmt"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/responses"
	"net/http"
	"os"
)

type customerSvc struct {
	customer repositories.CustomerRepo
}

func NewCustomerSvc(customer repositories.CustomerRepo) CustomerSvc {
	return &customerSvc{
		customer: customer,
	}
}

func (c *customerSvc) GetCustomers(ctx context.Context) *responses.Response {
	customers, err := c.customer.FindCustomers(ctx)
	if err != nil {
		return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
	}
	var response []responses.Customer
	for _, val := range customers {
		response = append(response, responses.Customer{
			Id:        val.Id,
			FullName:  val.FullName,
			KtpNumber: val.KtpNumber,
			Email:     val.Email,
		})
	}
	return responses.SuccessResponseWithData(responses.M_OK, http.StatusOK, response)
}

func (c *customerSvc) GetDetailCustomer(ctx context.Context, id uint) *responses.Response {
	customer, err := c.customer.FindCustomerById(ctx, id)
	if err != nil {
		return responses.ErrorResponse(responses.M_NOT_FOUND, http.StatusNotFound, errors.New("customer not found"))
	}
	response := responses.CustomerDetail{
		Id:              customer.Id,
		FullName:        customer.FullName,
		KtpNumber:       customer.KtpNumber,
		Gender:          customer.Gender,
		DateOfBirth:     customer.DateOfBirth,
		Address:         customer.Address,
		PhoneNumber:     customer.PhoneNumber,
		Email:           customer.Email,
		Nationality:     customer.Nationality,
		AddressProvince: customer.Province.Name,
		KtpImage:        fmt.Sprintf("%s%s/%s", os.Getenv("APP_URL"), ":"+os.Getenv("PORT"), customer.KtpImage),
		SelfieImage:     fmt.Sprintf("%s%s/%s", os.Getenv("APP_URL"), ":"+os.Getenv("PORT"), customer.SelfieImage),
		CreatedAt:       customer.CreatedAt,
		UpdatedAt:       customer.UpdatedAt,
	}
	return responses.SuccessResponseWithData(responses.M_OK, http.StatusOK, response)
}

func (c *customerSvc) GetCustomerLoanApplications(ctx context.Context, id uint) *responses.Response {
	customer, err := c.customer.FindCustomerById(ctx, id)
	if err != nil {
		return responses.ErrorResponse(responses.M_NOT_FOUND, http.StatusNotFound, errors.New("customer not found"))
	}
	var customerLoanRequest []responses.CustomerLoanRequest
	for _, val := range customer.CustomerLoanRequest {
		customerLoanRequest = append(customerLoanRequest, responses.CustomerLoanRequest{
			Id:        val.Id,
			Amount:    val.Amount,
			Tenor:     val.Tenor,
			Status:    val.Status,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
		})
	}
	return responses.SuccessResponseWithData(responses.M_OK, http.StatusOK, customerLoanRequest)
}
