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

type HTTPError struct {
	Message string
	Code    string
}

// 400

type BadRequestError HTTPError

func NewBadRequest(msg string, code string) *BadRequestError {
	return &BadRequestError{
		Message: msg,
		Code:    code,
	}
}

func (e *BadRequestError) Error() string {
	return e.Message
}

// 401

type UnauthorizedError HTTPError

func NewUnauthorizedError(msg string, code string) *UnauthorizedError {
	return &UnauthorizedError{
		Message: msg,
		Code:    code,
	}
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

// 403

type ForbiddenError HTTPError

func NewForbiddenError(msg string, code string) *ForbiddenError {
	return &ForbiddenError{
		Message: msg,
		Code:    code,
	}
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

// 404

type NotFoundError HTTPError

func NewNotFoundError(msg string, code string) *NotFoundError {
	return &NotFoundError{
		Message: msg,
		Code:    code,
	}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

/// 5xx status codes
