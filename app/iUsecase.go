package app

import (
	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
)

type IUsecase interface {
	CreateClient(name string) (clientID uuid.UUID, err error)
	GetClient(clientID uuid.UUID) (client models.Client, err error)

	CreateAccount(clientID uuid.UUID, currencyID uint) (accountID uuid.UUID, err error)
	GetAccount(accountID uuid.UUID) (account models.Account, err error)

	CreateTransaction(transaction models.Transaction) (transactionID uuid.UUID, err error)
	GetTransactions(accountID uuid.UUID) (list []models.Transaction, err error)
}
