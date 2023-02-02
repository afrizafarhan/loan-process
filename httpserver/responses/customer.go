package responses

import "time"

type Customer struct {
	Id        uint   `json:"id"`
	FullName  string `json:"full_name"`
	KtpNumber string `json:"ktp_number"`
	Email     string `json:"email"`
}

type CustomerDetail struct {
	Id              uint      `json:"id"`
	FullName        string    `json:"full_name"`
	KtpNumber       string    `json:"ktp_number"`
	Gender          string    `json:"gender"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Address         string    `json:"address"`
	PhoneNumber     string    `json:"phone_number"`
	Email           string    `json:"email"`
	Nationality     string    `json:"nationality"`
	AddressProvince string    `json:"address_province"`
	KtpImage        string    `json:"ktp_image"`
	SelfieImage     string    `json:"selfie_image"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CustomerLoanRequest struct {
	Id        uint      `json:"id"`
	Amount    uint      `json:"amount"`
	Tenor     uint      `json:"tenor"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
