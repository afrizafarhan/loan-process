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

func (c *customerRepo) FindCustomerByEmail(ctx context.Context, email string) (*models.Customer, error) {
	customer := new(models.Customer)
	err := c.db.WithContext(ctx).Where("email = ?", email).Take(customer).Error
	return customer, err
}

func (c *customerRepo) FindCustomerById(ctx context.Context, id uint) (*models.Customer, error) {
	customer := new(models.Customer)
	err := c.db.WithContext(ctx).Where("id = ?", id).Take(customer).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}
