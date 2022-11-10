package api

import (
	"banksystem/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("authorization")

		if authorizationHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenSlice := strings.Split(authorizationHeader, "Bearer")

		token := strings.Trim(tokenSlice[1], " ")

		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		payload, err := tokenMaker.VerifyToken(token)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
			return
		}

		ctx.Set("user", payload)
	}
}
