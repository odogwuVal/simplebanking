package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/odogwuVal/simplebanking/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD NGA"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	acc, err := server.store.CreateAccount(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, acc)
}
