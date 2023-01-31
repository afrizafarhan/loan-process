package gorm

import (
	"context"
	"gorm.io/gorm"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/repositories/models"
)

type customerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) repositories.CustomerRepo {
	return &customerRepo{
		db: db,
	}
}

func (c *customerRepo) SaveCustomer(ctx context.Context, customer *models.Customer) error {
	return c.db.WithContext(ctx).Save(customer).Error
}

func (c *customerRepo) FindCustomerByKtpNumber(ctx context.Context, ktpNumber string) (*models.Customer, error) {
	customer := new(models.Customer)
	err := c.db.WithContext(ctx).Where("ktp_number = ?", ktpNumber).Take(customer).Error
	return customer, err
}
