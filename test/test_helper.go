package test

import (
	"context"
	"gorm.io/gorm"
	gorm2 "loan_process/httpserver/repositories/gorm"
	"loan_process/httpserver/repositories/models"
	"time"
)

func createCustomerWithCustomerLoanRequest(db *gorm.DB, status string) (*models.Customer, *models.CustomerLoanRequest) {
	dateOfBirth, _ := time.Parse("2006-01-01", "2001-01-01")
	customer := models.Customer{
		FullName:    "Farhan",
		KtpNumber:   "1234567890123456",
		Gender:      "male",
		DateOfBirth: dateOfBirth,
		PhoneNumber: "081234567890",
		Email:       "farhan@gmail.com",
		Nationality: "indonesia",
		ProvinceId:  1,
		KtpImage:    "ktp/image.png",
		SelfieImage: "selfie/image.png",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	customerRepo := gorm2.NewCustomerRepo(db)
	customerLoanRequestRepo := gorm2.NewCustomerLoanRequestRepo(db)
	customerRepo.SaveCustomer(context.Background(), &customer)
	loanRequest := models.CustomerLoanRequest{
		CustomerId: customer.Id,
		Amount:     1000000,
		Tenor:      3,
		Status:     status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	customerLoanRequestRepo.SaveLoanRequest(context.Background(), &loanRequest)
	return &customer, &loanRequest
}
