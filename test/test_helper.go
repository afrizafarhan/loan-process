package test

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"loan_process/httpserver"
	"loan_process/httpserver/controllers"
	"loan_process/httpserver/middlewares"
	gorm2 "loan_process/httpserver/repositories/gorm"
	"loan_process/httpserver/repositories/models"
	"loan_process/httpserver/services"
	"time"
)

func createCustomerWithCustomerLoanRequest(db *gorm.DB, status string) (*models.Customer, *models.CustomerLoanRequest) {
	dateOfBirth, _ := time.Parse("2006-01-01", "2001-01-01")
	customer := models.Customer{
		FullName:    "Farhan",
		KtpNumber:   "1234567890123456",
		Gender:      "male",
		DateOfBirth: dateOfBirth,
		PhoneNumber: "081234567890",
		Email:       "farhan@gmail.com",
		Nationality: "indonesia",
		ProvinceId:  1,
		KtpImage:    "ktp/image.png",
		SelfieImage: "selfie/image.png",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	customerRepo := gorm2.NewCustomerRepo(db)
	customerLoanRequestRepo := gorm2.NewCustomerLoanRequestRepo(db)
	customerRepo.SaveCustomer(context.Background(), &customer)
	loanRequest := models.CustomerLoanRequest{
		CustomerId: customer.Id,
		Amount:     1000000,
		Tenor:      3,
		Status:     status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	customerLoanRequestRepo.SaveLoanRequest(context.Background(), &loanRequest)
	return &customer, &loanRequest
}

func createInstallment(db *gorm.DB, customer *models.Customer, customerLoanRequest *models.CustomerLoanRequest) *[]models.PaymentInstallment {
	var instalments []models.PaymentInstallment
	installmentAmount := customerLoanRequest.Amount / customerLoanRequest.Tenor
	year := time.Now().Year()
	month := time.Now().Month()
	day := 10
	for i := uint(1); i <= customerLoanRequest.Tenor; i++ {
		if month == 12 && i == 2 {
			month += 1
			year += 1
		}
		dueDate := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		instalments = append(instalments, models.PaymentInstallment{
			CustomerId:            customer.Id,
			CustomerLoanRequestId: customerLoanRequest.Id,
			Amount:                installmentAmount,
			Status:                "not_paid",
			DueDate:               dueDate,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		})
	}
	paymentInstallment := gorm2.NewPaymentInstallmentRepo(db)
	paymentInstallment.SavePaymentInstalment(context.Background(), &instalments)
	return &instalments
}

func GetCustomerById(db *gorm.DB, id uint) *models.Customer {
	customerRepo := gorm2.NewCustomerRepo(db)
	customer, _ := customerRepo.FindCustomerById(context.Background(), id)
	return customer
}

func createDailyLoanRequest(db *gorm.DB, request uint) {
	dailyLoan := models.DailyLoanRequest{
		CurrentDateRequest: time.Now(),
		Request:            request,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	dailyLoanRepo := gorm2.NewDailyLoanRequestRepo(db)
	dailyLoanRepo.SaveDailyLoanRequest(context.Background(), &dailyLoan)
}

func setupApp(engine *gin.Engine, db *gorm.DB) {
	//repo
	customer := gorm2.NewCustomerRepo(db)
	province := gorm2.NewProvinceRepo(db)
	loanRequest := gorm2.NewCustomerLoanRequestRepo(db)
	dailyLoan := gorm2.NewDailyLoanRequestRepo(db)
	paymentInstalment := gorm2.NewPaymentInstallmentRepo(db)
	//service
	service := services.NewLoanApplicationSvc(customer, province, loanRequest, dailyLoan, paymentInstalment)
	dailyLoanSvc := services.NewDailyLoanRequestSvc(dailyLoan)
	customerSvc := services.NewCustomerSvc(customer)
	paymentInstalmentSvc := services.NewPaymentInstallmentSvc(paymentInstalment)
	//controller
	controller := controllers.NewLoanApplicationController(service)
	customerController := controllers.NewCustomerController(customerSvc)
	paymentInstallmentController := controllers.NewPaymentInstallmentController(paymentInstalmentSvc)
	//middleware
	middleware := middlewares.NewCheckDailyRequestMiddleware(dailyLoanSvc)

	router := httpserver.NewRouter(engine, controller, customerController, paymentInstallmentController, middleware)
	router.SetRouter()
}

func truncateCustomer(db *gorm.DB) {
	db.Exec("TRUNCATE customers CASCADE")
}

func truncateDailyLoan(db *gorm.DB) {
	db.Exec("TRUNCATE daily_loan_requests CASCADE")
}

func formatDate(format string, date string) string {
	time, _ := time.Parse(format, date)
	return time.String()
}
