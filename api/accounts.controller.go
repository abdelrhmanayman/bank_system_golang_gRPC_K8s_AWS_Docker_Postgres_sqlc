package api

import (
	db "banksystem/db/sqlc"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func (server *Server) CreateAccountController(ctx *gin.Context) {
	var accountReq createAccountRequest

	err := ctx.ShouldBindJSON(&accountReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	accountArg := db.CreateAccountParams{
		Owner:    accountReq.Owner,
		Currency: accountReq.Currency,
	}

	account, err := server.store.CreateAccount(ctx, accountArg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountReq struct {
	ID int64 `uri:"ID" binding:"required,min=0"`
}

func (server *Server) GetAccountController(ctx *gin.Context) {
	var getAccReq getAccountReq

	err := ctx.ShouldBindUri(&getAccReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, getAccReq.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsReq struct {
	PageSize   int32 `form:"page_size" binding:"required,min=5,max=10"`
	PageNumber int32 `form:"page_number" binding:"required,min=1"`
}

func (server *Server) ListAccountsController(ctx *gin.Context) {
	var listAccReq listAccountsReq

	err := ctx.ShouldBindQuery(&listAccReq)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  listAccReq.PageSize,
		Offset: (listAccReq.PageNumber - 1) * listAccReq.PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
