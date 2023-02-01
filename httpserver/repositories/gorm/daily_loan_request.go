package gorm

import (
	"context"
	"gorm.io/gorm"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/repositories/models"
	"time"
)

type dailyLoanRequestRepo struct {
	db *gorm.DB
}

func NewDailyLoanRequestRepo(db *gorm.DB) repositories.DailyLoanRequestRepo {
	return &dailyLoanRequestRepo{db: db}
}

func (d *dailyLoanRequestRepo) SaveDailyLoanRequest(ctx context.Context, dailyLoan *models.DailyLoanRequest) error {
	return d.db.WithContext(ctx).Save(dailyLoan).Error
}

func (d *dailyLoanRequestRepo) FindDailyLoanRequestByDate(ctx context.Context, date time.Time) (*models.DailyLoanRequest, error) {
	dailyLoan := new(models.DailyLoanRequest)
	err := d.db.WithContext(ctx).Where("current_date_request = ?", date).Take(dailyLoan).Error

	if err != nil {
		return nil, err
	}
	return dailyLoan, nil
}

func (d *dailyLoanRequestRepo) UpdateDailyLoanRequestById(ctx context.Context, request *models.DailyLoanRequest) error {
	return d.db.WithContext(ctx).Model(request).Updates(*request).Error
}
