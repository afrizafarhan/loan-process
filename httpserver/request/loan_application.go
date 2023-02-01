package request

import "mime/multipart"

type CreateLoanApplication struct {
	FullName        string                `form:"full_name" validate:"required"`
	KtpNumber       string                `form:"ktp_number" validate:"required,numeric,min=16"`
	Gender          string                `form:"gender" validate:"required,oneof=male female"`
	DateOfBirth     string                `form:"date_of_birth" validate:"required"`
	Address         string                `form:"address" validate:"required"`
	PhoneNumber     string                `form:"phone_number" validate:"required,numeric,min=13"`
	Email           string                `form:"email" validate:"required,email"`
	Nationality     string                `form:"nationality" validate:"required,oneof=indonesia"`
	AddressProvince string                `form:"address_province" validate:"required"`
	KtpImage        *multipart.FileHeader `form:"ktp_image" validate:"required"`
	SelfieImage     *multipart.FileHeader `form:"selfie_image" validate:"required"`
	LoanAmount      uint                  `form:"loan_amount" validate:"required"`
	Tenor           uint                  `form:"tenor" validate:"required"`
}
