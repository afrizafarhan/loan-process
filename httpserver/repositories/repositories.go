package repositories

import (
	"context"
	"loan_process/httpserver/repositories/models"
	"time"
)

type CustomerRepo interface {
	SaveCustomer(ctx context.Context, customer *models.Customer) error
	FindCustomerByKtpNumber(ctx context.Context, ktpNumber string) (*models.Customer, error)
	FindCustomerByEmail(ctx context.Context, email string) (*models.Customer, error)
	FindCustomerById(ctx context.Context, id uint) (*models.Customer, error)
}

type CustomerLoanRequestRepo interface {
	SaveLoanRequest(ctx context.Context, request *models.CustomerLoanRequest) error
	FindAcceptedLoanRequestByCustomer(ctx context.Context, customerId uint) (*models.CustomerLoanRequest, error)
	FindAll(ctx context.Context) ([]models.CustomerLoanRequest, error)
}

type ProvinceRepo interface {
	FindProvinceByName(ctx context.Context, name string) (*models.Province, error)
}

type DailyLoanRequestRepo interface {
	SaveDailyLoanRequest(ctx context.Context, dailyLoan *models.DailyLoanRequest) error
	FindDailyLoanRequestByDate(ctx context.Context, date time.Time) (*models.DailyLoanRequest, error)
	UpdateDailyLoanRequestById(ctx context.Context, request *models.DailyLoanRequest) error
}

type PaymentInstallmentRepo interface {
	SavePaymentInstalment(ctx context.Context, installment *[]models.PaymentInstallment) error
	FindInstallmentByCustomerId(ctx context.Context, customerId uint) ([]models.PaymentInstallment, error)
}
