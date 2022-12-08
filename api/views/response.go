package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, data any) {
	ReturnWithCode(ctx, http.StatusOK, data)
}

func ReturnWithCode(ctx *gin.Context, code int, data any) {
	Response(ctx, code, "", data)
}

func Error(ctx *gin.Context, code int, err error) {
	ctx.Abort()
	Response(ctx, code, err.Error(), nil)
}

func Response(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(http.StatusOK, JsonResult{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

type JsonResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
