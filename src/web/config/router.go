package web

import (
	"algo/src/shared"

	"github.com/gin-gonic/gin"
)

type Router struct{}

func NewRouter(
	engine *gin.Engine,
) *Router {
	engine.GET("/", func(ctx *gin.Context) {
		shared.SuccessResp(ctx, "OK")
	})
	return new(Router)
}
