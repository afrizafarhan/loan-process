package httpserver

import (
	"github.com/gin-gonic/gin"
	"loan_process/httpserver/controllers"
)

type Router struct {
	router          *gin.Engine
	loanApplication controllers.LoanApplicationController
}

func NewRouter(engine *gin.Engine, loanApplication controllers.LoanApplicationController) *Router {
	return &Router{
		router:          engine,
		loanApplication: loanApplication,
	}
}

func (r *Router) SetRouter() *Router {
	r.router.Static("/resources/", "./resources")
	r.router.Use(cors)
	r.router.POST("/v1/loan-applications", r.loanApplication.Create)
	r.router.GET("/v1/loan-applications", r.loanApplication.GetLoanApplications)
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
