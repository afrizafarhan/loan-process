package responses

import "time"

type PaymentInstallmentResponse struct {
	Id        uint      `json:"id"`
	Amount    uint      `json:"amount"`
	DueDate   time.Time `json:"due_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
