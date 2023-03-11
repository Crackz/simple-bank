package api

import (
	"fmt"
	"net/http"

	db "github.com/crackz/simple-bank/db/sqlc"
	"github.com/crackz/simple-bank/token"
	"github.com/gin-gonic/gin"
)

type createTransferDto struct {
	ToAccountID   int64  `json:"toAccountID" binding:"required,min=1"`
	FromAccountID int64  `json:"fromAccountID" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=0,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var createDto createTransferDto

	if err := ctx.ShouldBindJSON(&createDto); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	fromAccount, err := server.checkAccountExist(ctx, createDto.FromAccountID)
	if err != nil {
		return
	}

	toAccount, err := server.checkAccountExist(ctx, createDto.ToAccountID)
	if err != nil {
		return
	}

	if !isValidAccountCurrency(fromAccount, createDto.Currency) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("from account doesn't support currency: %v", fromAccount.Currency)))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if !isOwner(fromAccount, authPayload.Username) {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("account id: %v doesn't belong to current user", createDto.FromAccountID)))
		return
	}

	if !isValidAccountCurrency(toAccount, createDto.Currency) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("to account doesn't support currency: %v", fromAccount.Currency)))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        createDto.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}
