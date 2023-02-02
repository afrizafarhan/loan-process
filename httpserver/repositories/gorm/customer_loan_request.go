package gorm

import (
	"context"
	"gorm.io/gorm"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/repositories/models"
)

type customerLoanRequestRepo struct {
	db *gorm.DB
}

func NewCustomerLoanRequestRepo(db *gorm.DB) repositories.CustomerLoanRequestRepo {
	return &customerLoanRequestRepo{db: db}
}

func (c *customerLoanRequestRepo) SaveLoanRequest(ctx context.Context, request *models.CustomerLoanRequest) error {
	return c.db.WithContext(ctx).Save(request).Error
}

func (c *customerLoanRequestRepo) FindAcceptedLoanRequestByCustomer(ctx context.Context, customerId uint) (*models.CustomerLoanRequest, error) {
	loanRequest := new(models.CustomerLoanRequest)
	err := c.db.WithContext(ctx).Where("customer_id = ?", customerId).Where("status = ?", "accepted").Take(loanRequest).Error
	if err != nil {
		return nil, err
	}
	return loanRequest, nil
}

func (c *customerLoanRequestRepo) FindAll(ctx context.Context) ([]models.CustomerLoanRequest, error) {
	var customerLoanRequests []models.CustomerLoanRequest

	if err := c.db.WithContext(ctx).Preload("Customer").Find(&customerLoanRequests).Error; err != nil {
		return customerLoanRequests, err
	}
	return customerLoanRequests, nil
}
