package repositories

import (
	"context"
	"loan_process/httpserver/repositories/models"
	"time"
)

type CustomerRepo interface {
	SaveCustomer(ctx context.Context, customer *models.Customer) error
	FindCustomerByKtpNumber(ctx context.Context, ktpNumber string) (*models.Customer, error)
}

type CustomerLoanRequest interface {
	SaveLoanRequest(ctx context.Context, request *models.CustomerLoanRequest) error
	FindAcceptedLoanRequestByCustomer(ctx context.Context, customerId uint) (*models.CustomerLoanRequest, error)
	FindLoanRequest(ctx context.Context) ([]models.CustomerLoanRequest, error)
}

type Province interface {
	FindProvinceByName(ctx context.Context, name string) (*models.Province, error)
}

type DailyLoanRequest interface {
	SaveDailyLoanRequest(ctx context.Context, dailyLoan *models.DailyLoanRequest) error
	FindDailyLoanRequestByDate(ctx context.Context, date time.Time) (*models.DailyLoanRequest, error)
	UpdateDailyLoanRequestById(ctx context.Context, request *models.DailyLoanRequest) error
}
