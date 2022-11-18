package app

import (
	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
)

type IRepository interface {
	CreateClient(client *models.Client) error
	GetClient(clientID uuid.UUID) (client models.Client, err error)

	CreateAccount(clientID uuid.UUID, account *models.Account) error
	GetAccountByID(accountID uuid.UUID) (account models.Account, err error)

	CreateTransaction(transaction *models.Transaction) error
	GetTransactionsByAccount(accountID uuid.UUID) (transactions []models.Transaction, err error)

	Close() error
}
