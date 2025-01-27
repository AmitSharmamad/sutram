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

func InternalServerErrorResp(ctx *gin.Context, msg string, code string) {
	ctx.IndentedJSON(http.StatusBadRequest, Error(msg, code))
}
