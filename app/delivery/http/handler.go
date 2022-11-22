package http

import (
	"net/http"

	"github.com/IrDeTen/money_processing_service.git/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	uc        app.IUsecase
	converter converter
}

func NewHandler(uc app.IUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

// CreateClient godoc
// @Summary 		Creating new client
// @Description  	Create new client with specified name
// @Accept       	json
// @Produce      	json
// @Param 			input body newClient true "Name for the client"
// @Success 		200 {string} string "New client UUID" 
// @Failure 		400 {object} errResponse
// @Failure 		500 {object} errResponse
// @Router 			/client [post]
func (h *Handler) CreateClient(c *gin.Context) {
	var client newClient

	if err := c.BindJSON(&client); err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}

	mClient, err := h.converter.ClientToModel(client)
	if err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.uc.CreateClient(mClient)
	if err != nil {
		h.errResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"client_id": id})
}

// GetClient godoc
// @Summary 		Retrieves client based on given ID 
// @Produce      	json
// @Param 			id path string true "Client UUID"
// @Success 		200 {object} outClient  
// @Failure 		400 {object} errResponse
// @Failure 		500 {object} errResponse
// @Router 			/client/{id} [get]
func (h *Handler) GetClient(c *gin.Context) {
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
	oClient := h.converter.ClientFromModel(mClient, accounts)
	c.JSON(http.StatusOK, map[string]interface{}{"client": oClient})

}

// CreateAccount godoc
// @Summary 		Creating new account
// @Description  	Create account with the specified currency for the client
// @Accept       	json
// @Produce      	json
// @Param 			input body newAccount true "Account Data"
// @Success 		200 {string} string "New account UUID" 
// @Failure 		400 {object} errResponse
// @Failure 		500 {object} errResponse
// @Router 			/account [post]
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

// GetAccount godoc
// @Summary 		Retrieves account based on given ID 
// @Produce      	json
// @Param 			id path string true "Account UUID"
// @Success 		200 {object} outAccount 
// @Failure 		400 {object} errResponse
// @Failure 		500 {object} errResponse
// @Router 			/account/{id} [get]
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
	outAcc := h.converter.AccountFromModel(account)
	c.JSON(http.StatusOK, map[string]interface{}{"account": outAcc})
}

// CreateAccount godoc
// @Summary 		Creating transaction
// @Description  	Create transaction based on transaction type, account IDs and transaction amount
// @Accept       	json
// @Produce      	json
// @Param 			input body newTransaction true "Transaction data"
// @Success 		200 {string} string "New transaction UUID" 
// @Failure 		400 {object} errResponse
// @Failure 		500 {object} errResponse
// @Router 			/transaction [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	var t newTransaction
	if err := c.BindJSON(&t); err != nil {
		h.errResponse(c, http.StatusBadRequest, err)
		return
	}
	transaction, err := h.converter.NewTransactionToModel(t)
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

// GetTransactions godoc
// @Summary 		Retrieves transactions list based on given account ID 
// @Produce      	json
// @Param 			id path string true "Account UUID"
// @Success 		200 {array} outTransaction 
// @Failure 		400 {object} errResponse
// @Failure 		500 {object} errResponse
// @Router 			/transaction/{id} [get]
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
		list = append(list, h.converter.TransactionFromModel(val))
	}
	c.JSON(http.StatusOK, map[string]interface{}{"transactions": list})
}
