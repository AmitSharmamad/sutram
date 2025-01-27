package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DataResp struct {
	Result bool `json:"result"`
	Data   any  `json:"data"`
}

func Data(data interface{}) *DataResp {
	return &DataResp{
		Result: true,
		Data:   data,
	}
}

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

func SuccessResp(ctx *gin.Context, data interface{}) {
	ctx.IndentedJSON(http.StatusOK, Data(data))
}

func BadReqResp(ctx *gin.Context, msg string, code string) {
	ctx.IndentedJSON(http.StatusBadRequest, Error(msg, code))
}

func InternalServerErrorResp(ctx *gin.Context, msg string, code string) {
	ctx.IndentedJSON(http.StatusBadRequest, Error(msg, code))
}
