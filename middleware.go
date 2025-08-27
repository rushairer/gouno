package gouno

import (
	"net/http"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(GlobalConfig.WebServerConfig.RequestTimeout),
		timeout.WithResponse(
			func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, ErrRequestTimeoutResponse)
			},
		),
	)
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(
		func(ctx *gin.Context, err any) {
			ctx.JSON(http.StatusOK, ErrInternalServerErrorResponse)
		},
	)
}
