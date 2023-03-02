package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/crackz/simple-bank/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountDto struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EURO"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var createDto createAccountDto

	if err := ctx.ShouldBindJSON(&createDto); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    createDto.Owner,
		Currency: createDto.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountParam struct {
	AccountID int64 `uri:"accountID" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var params getAccountParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	account, err := server.checkAccountExist(ctx, params.AccountID)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsQuery struct {
	Page  int32 `form:"page" binding:"min=1"`
	Limit int32 `form:"limit" binding:"min=1,max=100"`
}

func (server *Server) getAccounts(ctx *gin.Context) {
	var query getAccountsQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 1
	}

	arg := db.ListAccountsParams{
		Limit:  query.Limit,
		Offset: (query.Page - 1) * query.Limit,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountParam struct {
	AccountID int64 `uri:"accountID" binding:"required,min=1"`
}

type updateAccountDto struct {
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var dto updateAccountDto
	var params updateAccountParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	foundAccount, err := server.checkAccountExist(ctx, params.AccountID)
	if err != nil {
		return
	}

	arg := db.UpdateAccountParams{
		ID:      foundAccount.ID,
		Balance: dto.Balance,
	}

	updatedAccount, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, updatedAccount)
}

type deleteAccountParam struct {
	AccountID int64 `uri:"accountID" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var params deleteAccountParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	foundAccount, err := server.checkAccountExist(ctx, params.AccountID)
	if err != nil {
		return
	}

	err = server.store.DeleteAccount(ctx, foundAccount.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, foundAccount)
}

func (server *Server) checkAccountExist(ctx *gin.Context, id int64) (account db.Account, err error) {
	account, err = server.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("account not found")))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
	}

	return
}
