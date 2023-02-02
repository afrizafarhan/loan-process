package middlewares

import "github.com/gin-gonic/gin"

type CheckDailyRequestMiddleware interface {
	CheckDailyRequest() gin.HandlerFunc
}
