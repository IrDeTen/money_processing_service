package http

import (
	"github.com/IrDeTen/money_processing_service.git/app"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc app.IUsecase
}

func NewHandler(uc app.IUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) CreateClient(c *gin.Context) {
}

func (h *Handler) GetClient(c *gin.Context) {
}

func (h *Handler) CreateAccount(c *gin.Context) {
}

func (h *Handler) GetAccount(c *gin.Context) {
}

func (h *Handler) CreateTransaction(c *gin.Context) {
}

func (h *Handler) GetTransactions(c *gin.Context) {
}
