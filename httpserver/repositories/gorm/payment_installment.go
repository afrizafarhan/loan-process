package gorm

import (
	"context"
	"gorm.io/gorm"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/repositories/models"
)

type paymentInstallmentRepo struct {
	db *gorm.DB
}

func NewPaymentInstallmentRepo(db *gorm.DB) repositories.PaymentInstallmentRepo {
	return &paymentInstallmentRepo{db: db}
}

func (p *paymentInstallmentRepo) SavePaymentInstalment(ctx context.Context, installment *[]models.PaymentInstallment) error {
	return p.db.WithContext(ctx).Create(installment).Error
}

func (p *paymentInstallmentRepo) FindInstallmentByCustomerId(ctx context.Context, customerId uint) ([]models.PaymentInstallment, error) {
	var installments []models.PaymentInstallment
	err := p.db.Where("customer_id = ?", customerId).Find(&installments).Error
	if err != nil {
		return installments, err
	}
	return installments, nil
}
