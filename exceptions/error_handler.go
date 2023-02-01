package exceptions

import (
	"github.com/gin-gonic/gin"
	"loan_process/httpserver/responses"
	"net/http"
)

func ValidationError(ctx *gin.Context, err error) {
	if err != nil {
		errResponse := responses.ErrorResponse(responses.M_UNPROCESSABLE_ENTITY, http.StatusUnprocessableEntity, err)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, errResponse)
	}
}
