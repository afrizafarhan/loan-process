package models

import "time"

type Customer struct {
	Id          uint
	FullName    string
	KtpNumber   string
	Gender      string
	DateOfBirth time.Time
	Address     string
	PhoneNumber string
	Email       string
	Nationality string
	ProvinceId  uint
	KtpImage    string
	SelfieImage string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Province    Province `gorm:"foreignKey:ProvinceId;references:Id"`
}
