package app

import (
	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
)

type IRepository interface {
	CreateClient(client models.Client) (id uuid.UUID, err error)
	GetClient(clientID uuid.UUID) (client models.Client, err error)

	CreateAccount(account models.Account) (id uuid.UUID, err error)
	GetAccountByID(accountID uuid.UUID) (account models.Account, err error)
	GetAccountsByClientID(clientID uuid.UUID) (accounts []models.Account, err error)

	CreateTransaction(transaction models.Transaction, accounts ...models.Account) (id uuid.UUID, err error)
	GetTransactionsByAccount(accountID uuid.UUID) (transactions []models.Transaction, err error)

	Close() error
}
