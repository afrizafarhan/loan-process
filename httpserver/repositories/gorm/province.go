package gorm

import (
	"context"
	"gorm.io/gorm"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/repositories/models"
)

type provinceRepo struct {
	db *gorm.DB
}

func NewProvinceRepo(db *gorm.DB) repositories.Province {
	return &provinceRepo{db: db}
}

func (p *provinceRepo) FindProvinceByName(ctx context.Context, name string) (*models.Province, error) {
	var province models.Province
	err := p.db.WithContext(ctx).Where("name = ?", name).Find(province).Error
	if err != nil {
		return &province, err
	}
	return &province, nil
}
