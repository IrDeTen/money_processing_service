package http

import (
	"github.com/IrDeTen/money_processing_service.git/app"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc app.IUsecase) {
	h := NewHandler(uc)

	//TODO change endpoint name
	apiEndpoints := router.Group("/api")
	{
		apiEndpoints.POST("", h.CreateClient)
		apiEndpoints.GET("", h.GetClient)

		apiEndpoints.POST("", h.CreateAccount)
		apiEndpoints.GET("", h.GetAccount)

		apiEndpoints.POST("", h.CreateTransaction)
		apiEndpoints.GET("", h.GetTransactions)
	}
}
