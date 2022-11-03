package api

import (
	db "banksystem/db/sqlc"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateTransferController(ctx *gin.Context) {
	var transferReq createTransferRequest

	err := ctx.ShouldBindJSON(&transferReq)

	if err != nil {
		fmt.Printf("%#v", err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	transferArg := db.CreateTransferParams{
		FromAccount: transferReq.FromAccountID,
		ToAccount:   transferReq.ToAccountID,
		Amount:      transferReq.Amount,
	}

	transfer, err := server.store.CreateTransfer(ctx, transferArg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}
