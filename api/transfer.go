package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/thaian1234/simplebank/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !server.validServer(ctx, req.FromAccountID, req.Currency) {
		ctx.JSON(http.StatusNotAcceptable, msgResponse("Invalid currency"))
		return
	}
	if !server.validServer(ctx, req.ToAccountID, req.Currency) {
		ctx.JSON(http.StatusNotAcceptable, msgResponse("Invalid currency"))
		return
	}
	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	transfer, err := server.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) validServer(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		return false
	}
	if account.Currency != currency {
		return false
	}
	return true
}
