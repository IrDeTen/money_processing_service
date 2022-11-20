package http

import (
	"fmt"
	"net/http"

	"github.com/IrDeTen/money_processing_service.git/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	var client newClient

	if err := c.BindJSON(&client); err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.uc.CreateClient(client.ToModel())
	if err != nil {
		h.errResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"client_id": id.String()})
}

func (h *Handler) GetClient(c *gin.Context) {
	var oClient outClient
	clientID := c.Param("client_id")
	id, err := uuid.Parse(clientID)
	if err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}

	mClient, accounts, err := h.uc.GetClient(id)
	if err != nil {
		h.errResponse(c, http.StatusInternalServerError, err)
		return
	}
	fmt.Println(mClient)
	fmt.Println(accounts)
	oClient.FromModel(mClient, accounts)
	c.JSON(http.StatusOK, map[string]interface{}{"client": oClient})

}

func (h *Handler) CreateAccount(c *gin.Context) {
	var newAcc newAccount

	if err := c.BindJSON(&newAcc); err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}
	clientID, err := uuid.Parse(newAcc.ClientID)
	if err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.uc.CreateAccount(clientID, newAcc.CurrencyID)
	if err != nil {
		h.errResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"account_id": id.String()})
}

func (h *Handler) GetAccount(c *gin.Context) {
	accountID := c.Param("account_id")
	id, err := uuid.Parse(accountID)
	if err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}

	account, err := h.uc.GetAccount(id)
	if err != nil {
		h.errResponse(c, http.StatusInternalServerError, err)
		return
	}
	outAcc := outAccount.FromModel(outAccount{}, account)
	c.JSON(http.StatusOK, map[string]interface{}{"account": outAcc})
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var ta newTransaction
	if err := c.BindJSON(&ta); err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}
	transaction, err := ta.ToModel()
	if err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.uc.CreateTransaction(transaction)
	if err != nil {
		h.errResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"transaction_id": id.String()})
}

func (h *Handler) GetTransactions(c *gin.Context) {
	accountID := c.Param("account_id")
	id, err := uuid.Parse(accountID)
	if err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}

	transactions, err := h.uc.GetTransactions(id)
	if err != nil {
		h.errResponse(c, http.StatusInternalServerError, err)
		return
	}

	list := make([]outTransaction, 0)
	for _, val := range transactions {
		list = append(list, outTransaction.FromModel(outTransaction{}, val))
	}
	c.JSON(http.StatusOK, map[string]interface{}{"transactions": list})
}
