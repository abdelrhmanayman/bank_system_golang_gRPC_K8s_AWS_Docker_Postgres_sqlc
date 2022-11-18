package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewTokenResponse struct {
	NewAccessToken          string    `json:"new_access_token"`
	NewAccessTokenExpiresAt time.Time `json:"new_access_token_expires_at"`
}

func (server *Server) RenewTokenController(ctx *gin.Context) {
	var renewTokenReq RenewTokenRequest

	err := ctx.ShouldBindJSON(&renewTokenReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(renewTokenReq.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusForbidden, ErrorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, payload.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusForbidden, ErrorResponse(errors.New("bara ya m3aras")))
		return
	}

	newAccessToken, payload, err := server.tokenMaker.CreateToken(payload.Username, server.appConfig.TokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, RenewTokenResponse{
		NewAccessToken:          newAccessToken,
		NewAccessTokenExpiresAt: payload.ExpiredAt,
	})
}
