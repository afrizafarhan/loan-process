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

func TestCustomer_GetSuccess(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)
	customer, _ := createCustomerWithCustomerLoanRequest(db, "accepted")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/customers", nil)
	req.Header.Add("Content-type", "application/json")
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	var loanApplications = responseBody["data"].([]interface{})
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	loanApplications1 := loanApplications[0].(map[string]interface{})
	assert.Equal(t, customer.FullName, loanApplications1["full_name"])
	assert.Equal(t, customer.KtpNumber, loanApplications1["ktp_number"])
	assert.Equal(t, customer.Email, loanApplications1["email"])
	assert.Equal(t, int(customer.Id), int(loanApplications1["id"].(float64)))
}

func TestCustomer_GetSuccessNoData(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/customers", nil)
	req.Header.Add("Content-type", "application/json")
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Nil(t, responseBody["data"])
}

func TestCustomer_GetDetailCustomerNotFound(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/customers/999", nil)
	req.Header.Add("Content-type", "application/json")
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "NOT_FOUND", responseBody["status"])
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "customer not found", responseBody["error"])
}

func TestCustomer_GetDetailCustomerSuccess(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)
	customer, _ := createCustomerWithCustomerLoanRequest(db, "accepted")
	customerDetail := GetCustomerById(db, customer.Id)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/customers/"+strconv.Itoa(int(customer.Id)), nil)
	req.Header.Add("Content-type", "application/json")
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	var loanApplications = responseBody["data"].(map[string]interface{})

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, int(customerDetail.Id), int(loanApplications["id"].(float64)))
	assert.Equal(t, customerDetail.FullName, loanApplications["full_name"])
	assert.Equal(t, customerDetail.KtpNumber, loanApplications["ktp_number"])
	assert.Equal(t, customerDetail.Gender, loanApplications["gender"])
	assert.Equal(t, formatDate("2006-01-01", customerDetail.DateOfBirth.String()),
		formatDate("2006-01-01", loanApplications["date_of_birth"].(string)),
	)
	assert.Equal(t, customerDetail.PhoneNumber, loanApplications["phone_number"])
	assert.Equal(t, customerDetail.Email, loanApplications["email"])
	assert.Equal(t, customerDetail.Nationality, loanApplications["nationality"])
	assert.Equal(t, customerDetail.Province.Name, loanApplications["address_province"])
	assert.Equal(t, formatDate("2006-01-01 00:00:00", customerDetail.CreatedAt.String()),
		formatDate("2006-01-01 00:00:00", loanApplications["created_at"].(string)),
	)
	assert.Equal(t, formatDate("2006-01-01 00:00:00", customerDetail.UpdatedAt.String()),
		formatDate("2006-01-01 00:00:00", loanApplications["updated_at"].(string)),
	)
}
