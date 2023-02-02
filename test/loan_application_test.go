package test

import (
	"bytes"
	"encoding/json"
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
	"strconv"
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

func TestLoanApplication_GetLoanApplicationSuccess(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/loan-applications", nil)
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	var loanApplications = responseBody["data"].([]interface{})
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	loanApplications1 := loanApplications[0].(map[string]interface{})
	assert.Equal(t, "Farhan", loanApplications1["full_name"])
	assert.Equal(t, "farhan@gmail.com", loanApplications1["email"])
	assert.Equal(t, "1234567890123456", loanApplications1["ktp_number"])
	assert.Equal(t, 1000000, int(loanApplications1["loan_amount"].(float64)))
	assert.Equal(t, 3, int(loanApplications1["tenor"].(float64)))
}

func TestLoanApplication_GetLoanReApplyFailedCauseDailyLimit(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications", nil)
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	var loanApplications = responseBody["data"].([]interface{})
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	loanApplications1 := loanApplications[0].(map[string]interface{})
	assert.Equal(t, "Farhan", loanApplications1["full_name"])
	assert.Equal(t, "farhan@gmail.com", loanApplications1["email"])
	assert.Equal(t, "1234567890123456", loanApplications1["ktp_number"])
	assert.Equal(t, 1000000, int(loanApplications1["loan_amount"].(float64)))
	assert.Equal(t, 3, int(loanApplications1["tenor"].(float64)))
}

func TestLoanApplication_PostLoanReApplyFailedCauseAlreadyAcceptedLoan(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)
	customer, _ := createCustomerWithCustomerLoanRequest(db, "accepted")

	body := strings.NewReader(`{"loan_amount":1000000,"tenor":6}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/loan-applications/"+strconv.Itoa(int(customer.Id))+"/reapply", body)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "BAD_REQUEST", responseBody["status"])
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "the customer already have accepted loan", responseBody["error"])
}

func TestLoanApplication_PostLoanReApplySuccess(t *testing.T) {
	router := gin.Default()
	db, _ := config.ConnectPostgresGORMTest()
	setupApp(router, db)
	truncateCustomer(db)
	customer, _ := createCustomerWithCustomerLoanRequest(db, "rejected")

	w := httptest.NewRecorder()
	body := strings.NewReader(`{"loan_amount":1000000,"tenor":6}`)
	req, _ := http.NewRequest("POST", "/v1/loan-applications/"+strconv.Itoa(int(customer.Id))+"/reapply", body)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	response, _ := io.ReadAll(w.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(response, &responseBody)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "Success reapply loan application", responseBody["message"])
}
