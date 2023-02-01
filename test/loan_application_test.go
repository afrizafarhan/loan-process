package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"io"
	"loan_process/config"
	"loan_process/httpserver"
	"loan_process/httpserver/controllers"
	repo "loan_process/httpserver/repositories/gorm"
	"loan_process/httpserver/services"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupApp(engine *gin.Engine, db *gorm.DB) {
	//repo
	customer := repo.NewCustomerRepo(db)
	province := repo.NewProvinceRepo(db)
	loanRequest := repo.NewCustomerLoanRequestRepo(db)
	dailyLoan := repo.NewDailyLoanRequestRepo(db)
	paymentInstalment := repo.NewPaymentInstallmentRepo(db)
	//service
	service := services.NewLoanApplicationSvc(customer, province, loanRequest, dailyLoan, paymentInstalment)
	controller := controllers.NewLoanApplicationController(service)
	router := httpserver.NewRouter(engine, controller)
	router.SetRouter()
}

func truncateCustomer(db *gorm.DB) {
	db.Exec("TRUNCATE customers CASCADE")
}

func truncateDailyLoan(db *gorm.DB) {
	db.Exec("TRUNCATE daily_loan_requests CASCADE")
}

func TestLoanApplication_BadRequestNotGivingPayload(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	w := httptest.NewRecorder()
	body := strings.NewReader(``)
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
	fmt.Println(w.Body)
}

func createPayload(writer *multipart.Writer, payload *map[string]string, incorrect bool) {
	pathFile := "images/test_image.png"
	if incorrect {
		pathFile = "txt/test.txt"
	}
	file, _ := os.Open(pathFile)
	defer file.Close()

	ktp, _ := writer.CreateFormFile("ktp_image", filepath.Base(pathFile))
	selfie, _ := writer.CreateFormFile("selfie_image", filepath.Base(pathFile))

	io.Copy(ktp, file)
	io.Copy(selfie, file)
	for key, val := range *payload {
		writer.WriteField(key, val)
	}
	writer.Close()
}

func TestLoanApplication_CreateValidateFailed(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	var payload = map[string]string{
		"full_name":        "Farhan",
		"ktp_number":       "1234567890",
		"gender":           "aaaa",
		"date_of_birth":    "2001-01-01",
		"address":          "Jl. Test",
		"phone_number":     "08123450",
		"email":            "farhan",
		"nationality":      "aaa",
		"address_province": "SUMATERA UTARA",
		"loan_amount":      "10",
		"tenor":            "3",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	createPayload(writer, &payload, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 422, w.Code)
	assert.Equal(t, "UNPROCESSABLE_ENTITY", responseBody["status"])
	assert.Equal(t, 422, int(responseBody["code"].(float64)))
	assert.NotEmpty(t, responseBody["error"])
}

func TestLoanApplication_CreateFailedCauseUnderAge(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	var payload = map[string]string{
		"full_name":        "Farhan",
		"ktp_number":       "1234567890123456",
		"gender":           "male",
		"date_of_birth":    "2021-01-01",
		"address":          "Jl. Test",
		"phone_number":     "08123456890",
		"email":            "farhan@gmail.com",
		"nationality":      "indonesia",
		"address_province": "SUMATERA UTARA",
		"loan_amount":      "1000000",
		"tenor":            "3",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	createPayload(writer, &payload, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 422, w.Code)
	assert.Equal(t, "UNPROCESSABLE_ENTITY", responseBody["status"])
	assert.Equal(t, 422, int(responseBody["code"].(float64)))
	assert.Equal(t, "age must between 17 and 80", responseBody["error"])
}

func TestLoanApplication_CreateFailedCauseProvinceNotFound(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	var payload = map[string]string{
		"full_name":        "Farhan",
		"ktp_number":       "1234567890123456",
		"gender":           "male",
		"date_of_birth":    "2001-01-01",
		"address":          "Jl. Test",
		"phone_number":     "08123456890",
		"email":            "farhan@gmail.com",
		"nationality":      "indonesia",
		"address_province": "KALIMANTAN UTARA",
		"loan_amount":      "1000000",
		"tenor":            "3",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	createPayload(writer, &payload, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 422, w.Code)
	assert.Equal(t, "UNPROCESSABLE_ENTITY", responseBody["status"])
	assert.Equal(t, 422, int(responseBody["code"].(float64)))
	assert.Equal(t, "the province not allowed to loan application request", responseBody["error"])
}

func TestLoanApplication_CreateFailedCauseTenorNotCorrect(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	var payload = map[string]string{
		"full_name":        "Farhan",
		"ktp_number":       "1234567890123456",
		"gender":           "male",
		"date_of_birth":    "2001-01-01",
		"address":          "Jl. Test",
		"phone_number":     "08123456890",
		"email":            "farhan@gmail.com",
		"nationality":      "indonesia",
		"address_province": "SUMATERA UTARA",
		"loan_amount":      "1000000",
		"tenor":            "1",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	createPayload(writer, &payload, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 422, w.Code)
	assert.Equal(t, "UNPROCESSABLE_ENTITY", responseBody["status"])
	assert.Equal(t, 422, int(responseBody["code"].(float64)))
	assert.Equal(t, "the tenor not in right value", responseBody["error"])
}
func TestLoanApplication_CreateFailedCauseImageExt(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	var payload = map[string]string{
		"full_name":        "Farhan",
		"ktp_number":       "1234567890123456",
		"gender":           "male",
		"date_of_birth":    "2001-01-01",
		"address":          "Jl. Test",
		"phone_number":     "08123456890",
		"email":            "farhan@gmail.com",
		"nationality":      "indonesia",
		"address_province": "SUMATERA UTARA",
		"loan_amount":      "1000000",
		"tenor":            "3",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	createPayload(writer, &payload, true)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 422, w.Code)
	assert.Equal(t, "UNPROCESSABLE_ENTITY", responseBody["status"])
	assert.Equal(t, 422, int(responseBody["code"].(float64)))
	assert.Equal(t, "ktp image extension not allowed", responseBody["error"])
}

func TestLoanApplication_CreateSuccess(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	var payload = map[string]string{
		"full_name":        "Farhan",
		"ktp_number":       "1234567890123456",
		"gender":           "male",
		"date_of_birth":    "2001-01-01",
		"address":          "Jl. Test",
		"phone_number":     "08123456890",
		"email":            "farhan@gmail.com",
		"nationality":      "indonesia",
		"address_province": "SUMATERA UTARA",
		"loan_amount":      "1000000",
		"tenor":            "3",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	createPayload(writer, &payload, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "CREATED", responseBody["status"])
	assert.Equal(t, 201, int(responseBody["code"].(float64)))
	assert.Equal(t, "Success create loan application", responseBody["message"])
}

func TestLoanApplication_CreateFailedKtpNumberAlreadyExists(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)

	var payload = map[string]string{
		"full_name":        "Farhan",
		"ktp_number":       "1234567890123456",
		"gender":           "male",
		"date_of_birth":    "2001-01-01",
		"address":          "Jl. Test",
		"phone_number":     "08123456890",
		"email":            "farhan@gmail.com",
		"nationality":      "indonesia",
		"address_province": "SUMATERA UTARA",
		"loan_amount":      "1000000",
		"tenor":            "3",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	createPayload(writer, &payload, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "BAD_REQUEST", responseBody["status"])
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "ktp number already exist", responseBody["error"])
}

func TestLoanApplication_CreateFailedEmailAlreadyExists(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)

	var payload = map[string]string{
		"full_name":        "Farhan",
		"ktp_number":       "1234567890123457",
		"gender":           "male",
		"date_of_birth":    "2001-01-01",
		"address":          "Jl. Test",
		"phone_number":     "08123456890",
		"email":            "farhan@gmail.com",
		"nationality":      "indonesia",
		"address_province": "SUMATERA UTARA",
		"loan_amount":      "1000000",
		"tenor":            "3",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	createPayload(writer, &payload, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "BAD_REQUEST", responseBody["status"])
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "email already exist", responseBody["error"])
}
