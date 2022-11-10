package api

import (
	db "banksystem/db/sqlc"
	"banksystem/util"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User           db.User   `json:"user"`
	Token          string    `json:"token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
}

var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
)

func (server *Server) LoginController(ctx *gin.Context) {
	var loginReq LoginRequest
	var user db.User

	err := ctx.ShouldBindJSON(&loginReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	user, err = server.store.GetUser(ctx, loginReq.Username)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			ctx.JSON(http.StatusForbidden, ErrorResponse(ErrInvalidCredentials))
			return
		}

		ctx.JSON(http.StatusForbidden, ErrorResponse(err))
		return
	}

	isPasswordValid := util.CompareHashedPasswords(loginReq.Password, user.HashedPwd)

	if !isPasswordValid {
		ctx.JSON(http.StatusForbidden, ErrorResponse(ErrInvalidCredentials))
		return
	}

	token, payload, err := server.tokenMaker.CreateToken(user.Username, server.appConfig.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		User: db.User{
			Username: user.Username,
			Email:    user.Email,
		},
		Token:          token,
		TokenExpiresAt: payload.ExpiredAt,
	})

}
