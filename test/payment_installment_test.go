package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"loan_process/config"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCustomer_GetInstallmentCustomerFailed(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	customer, _ := createCustomerWithCustomerLoanRequest(db, "rejected")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/payment-installments/"+strconv.Itoa(int(customer.Id)), nil)
	req.Header.Add("Content-type", "application/json")
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "NOT_FOUND", responseBody["status"])
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "customer accepted loan request not found", responseBody["error"])
}

func TestCustomer_GetInstallmentCustomerSuccess(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)
	customer, loanRequest := createCustomerWithCustomerLoanRequest(db, "accepted")
	installment := createInstallment(db, customer, loanRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/payment-installments/"+strconv.Itoa(int(customer.Id)), nil)
	req.Header.Add("Content-type", "application/json")
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	var customers = responseBody["data"].([]interface{})
	customer1 := customers[0].(map[string]interface{})
	assert.Equal(t, int((*installment)[0].Id), int(customer1["id"].(float64)))
	assert.Equal(t, int((*installment)[0].Amount), int(customer1["amount"].(float64)))
	assert.Equal(t, formatDate("2006-01-01", (*installment)[0].DueDate.String()), formatDate("2006-01-01", customer1["due_date"].(string)))
	assert.Equal(t, (*installment)[0].Status, customer1["status"])
	assert.Equal(t, formatDate("2006-01-01 00:00:00", (*installment)[0].CreatedAt.String()),
		formatDate("2006-01-01 00:00:00", customer1["created_at"].(string)),
	)
	assert.Equal(t, formatDate("2006-01-01 00:00:00", (*installment)[0].UpdatedAt.String()),
		formatDate("2006-01-01 00:00:00", customer1["updated_at"].(string)),
	)
}
