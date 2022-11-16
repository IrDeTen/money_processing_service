package http

import (
	"github.com/IrDeTen/money_processing_service.git/app"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc app.IUsecase) {
	h := NewHandler(uc)

	clientEndpoints := router.Group("/client")
	{
		clientEndpoints.POST("", h.CreateClient)
		clientEndpoints.GET(":clent_id", h.GetClient)
	}

	accountEndpoints := router.Group("/account")
	{
		accountEndpoints.POST("", h.CreateAccount)
		accountEndpoints.GET(":account_id", h.GetAccount)

	}

	transactionEndpoints := router.Group("/transaction")
	{
		transactionEndpoints.POST("", h.CreateTransaction)
		transactionEndpoints.GET(":account_id", h.GetTransactions)
	}
}
