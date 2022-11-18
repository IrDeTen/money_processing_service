package postgres

import (
	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
)

type Repository struct {
}

func NewRepo() *Repository {
	return &Repository{}
}

func (r *Repository) CreateClient(client *models.Client) error

func (r *Repository) GetClient(clientID uuid.UUID) (client models.Client, err error)

func (r *Repository) CreateAccount(clientID uuid.UUID, account *models.Account) error

func (r *Repository) GetAccountByID(accountID uuid.UUID) (account models.Account, err error)

func (r *Repository) UpdateBalance(account *models.Account) error

func (r *Repository) UpdateMultipleBalances(account ...*models.Account) error

func (r *Repository) CreateTransaction(transaction *models.Transaction) error

func (r *Repository) GetTransactionsByAccount(accountID uuid.UUID) (transactions []models.Transaction, err error)

func (r *Repository) Close() error {
	return nil
}
