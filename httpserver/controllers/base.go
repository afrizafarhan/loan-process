package controllers

import (
	"github.com/gin-gonic/gin"
	"loan_process/httpserver/responses"
)

func WriteJsonResponse(ctx *gin.Context, res *responses.Response) {
	ctx.JSON(res.Code, res)
}
