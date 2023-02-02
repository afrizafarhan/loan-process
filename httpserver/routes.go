package httpserver

import (
	"github.com/gin-gonic/gin"
	"loan_process/httpserver/controllers"
	"loan_process/httpserver/middlewares"
)

type Router struct {
	router          *gin.Engine
	loanApplication controllers.LoanApplicationController
	customer        controllers.CustomerController
	middleware      middlewares.CheckDailyRequestMiddleware
}

func NewRouter(engine *gin.Engine, loanApplication controllers.LoanApplicationController, customer controllers.CustomerController, middleware middlewares.CheckDailyRequestMiddleware) *Router {
	return &Router{
		router:          engine,
		loanApplication: loanApplication,
		customer:        customer,
		middleware:      middleware,
	}
}

func (r *Router) SetRouter() *Router {
	r.router.Static("/resources/", "./resources")
	r.router.Use(cors)

	r.router.POST("/v1/loan-applications", r.middleware.CheckDailyRequest(), r.loanApplication.Create)
	r.router.POST("/v1/loan-applications/:customerId/reapply", r.middleware.CheckDailyRequest(), r.loanApplication.ReapplyLoanApplication)
	r.router.GET("/v1/loan-applications", r.loanApplication.GetLoanApplications)

	r.router.GET("/v1/customers", r.customer.GetCustomers)
	r.router.GET("/v1/customers/:id", r.customer.GetDetailCustomer)
	r.router.GET("/v1/customers/:id/loan-applications", r.customer.GetCustomerLoanApplications)

	return r
}

func (r *Router) Start(port string) {
	r.SetRouter()
	r.router.Run(port)
}

func cors(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
