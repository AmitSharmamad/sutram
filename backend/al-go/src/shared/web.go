package shared

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/// 2xx status codes

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

func SuccessResp(ctx *gin.Context, data interface{}) {
	ctx.IndentedJSON(http.StatusOK, Data(data))
}

/// 4xx status codes

type BadRequestError struct {
	Message string
	Code    string
}

func NewBadRequest(msg string, code string) *BadRequestError {
	return &BadRequestError{
		Message: msg,
		Code:    code,
	}
}

func (e *BadRequestError) Error() string {
	return e.Message
}

/// 5xx status codes
