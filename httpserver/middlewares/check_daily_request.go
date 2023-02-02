package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"loan_process/httpserver/responses"
	"loan_process/httpserver/services"
	"net/http"
)

type checkDailyRequestMiddleware struct {
	svc services.DailyLoanRequestSvc
}

func NewCheckDailyRequestMiddleware(svc services.DailyLoanRequestSvc) CheckDailyRequestMiddleware {
	return &checkDailyRequestMiddleware{svc: svc}
}

func (d *checkDailyRequestMiddleware) CheckDailyRequest() gin.HandlerFunc {
	return func(context *gin.Context) {
		if d.svc.CheckDailyLoanRequest(context) < 50 {
			context.Next()
			return
		}
		response := responses.ErrorResponse(responses.M_BAD_REQUEST, http.StatusBadRequest, errors.New("the loan application daily limit exceeded"))
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
}
