package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/odogwuVal/simplebanking/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validateAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validateAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) bool {
	acc, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	if acc.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s & %s", accountID, currency, acc.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
