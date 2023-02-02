package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"loan_process/config"
	"loan_process/httpserver"
	"loan_process/httpserver/controllers"
	"loan_process/httpserver/middlewares"
	gorm2 "loan_process/httpserver/repositories/gorm"
	"loan_process/httpserver/services"
	"log"
	"os"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("cannot load your env")
	}
}

func main() {
	db, err := config.ConnectPostgresGORM()
	if err != nil {
		panic(err)
	}
	if err = os.Mkdir("./resources/ktp", os.ModePerm); err != nil {
		log.Println(err)
	}
	if err = os.Mkdir("./resources/selfie", os.ModePerm); err != nil {
		log.Println(err)
	}
	os.Setenv("APP_ENV", "production")

	router := gin.Default()
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

	app := httpserver.NewRouter(router, controller, customerController, paymentInstallmentController, middleware)
	PORT := os.Getenv("PORT")
	app.Start(":" + PORT)
}
