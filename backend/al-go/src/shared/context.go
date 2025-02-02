package shared

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ReqCtx struct {
	UserID string
}

type ReqParam string

const (
	UserID ReqParam = "UserID"
)

func GetReqCtx(
	gctx *gin.Context,
) context.Context {
	ctx := context.Background()
	userID, exists := gctx.Get("user_id")
	if !exists {
		panic("Invalid user context")
	}
	return context.WithValue(ctx, UserID, userID)
}

func FromReqCtx[T any](ctx context.Context, key ReqParam) T {
	return ctx.Value(key).(T)
}
