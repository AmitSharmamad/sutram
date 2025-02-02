package web

import (
	"fmt"

	"algo/src/shared"

	"github.com/gin-gonic/gin"
)

func GlobalErrorHandler(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Log the error (you might want to use a logging library)
			fmt.Println("Recovered from panic:")
			fmt.Println(err)

			switch err := err.(type) {

			case *shared.BadRequestError:
				BadReqResp(ctx, err.Message, err.Code)
				return

			case *shared.UnauthorizedError:
				UnauthorizedResp(ctx, err.Message, err.Code)
				return

			case *shared.ForbiddenError:
				ForbiddenResp(ctx, err.Message, err.Code)
				return

			default:
				InternalServerErrorResp(ctx, "Internal Server Error", "INTERNAL_SERVER_ERROR")
				return
			}

		}
	}()

	ctx.Next()
}
