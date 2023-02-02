package responses

type CustomerLoanRequestResponses struct {
	Id         uint   `json:"id"`
	FullName   string `json:"full_name"`
	KtpNumber  string `json:"ktp_number"`
	Email      string `json:"email"`
	LoanAmount uint   `json:"loan_amount"`
	Tenor      uint   `json:"tenor"`
	Status     string `json:"status"`
}
