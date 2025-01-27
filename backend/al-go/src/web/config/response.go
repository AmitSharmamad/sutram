package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrResp struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func Error(msg string, code string) *ErrResp {
	return &ErrResp{
		Result:  false,
		Message: msg,
		Code:    code,
	}
}

func BadReqResp(ctx *gin.Context, msg string, code string) {
	ctx.IndentedJSON(http.StatusBadRequest, Error(msg, code))
}

func UnauthorizedResp(ctx *gin.Context, msg string, code string) {
	ctx.IndentedJSON(http.StatusUnauthorized, nil)
}

func ForbiddenResp(ctx *gin.Context, msg string, code string) {
	ctx.IndentedJSON(http.StatusForbidden, nil)
}

func InternalServerErrorResp(ctx *gin.Context, msg string, code string) {
	ctx.IndentedJSON(http.StatusBadRequest, Error(msg, code))
}
