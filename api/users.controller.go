package api

import (
	db "banksystem/db/sqlc"
	"banksystem/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (server *Server) CreateUserController(ctx *gin.Context) {
	var createUserReq CreateUserRequest

	err := ctx.ShouldBindJSON(&createUserReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	userHashedPassword, err := util.HashPassword(createUserReq.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	userDbArgs := db.CreateUserParams{
		Username:  createUserReq.Username,
		HashedPwd: userHashedPassword,
		Email:     createUserReq.Email,
	}

	user, err := server.store.CreateUser(ctx, userDbArgs)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}
