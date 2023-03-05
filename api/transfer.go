package api

import (
	"net/http"

	db "github.com/crackz/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferDto struct {
	ToAccountID   int64 `json:"toAccountID" binding:"required,min=1"`
	FromAccountID int64 `json:"fromAccountID" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,min=0,gt=0"`
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

	// TODO: Convert Amount From FromAccount's Currency To ToAccount's Currency

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
