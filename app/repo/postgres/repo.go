package postgres

import (
	"database/sql"

	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (r *Repository) CreateClient(client *models.Client) error {
	return nil
}

func (r *Repository) GetClient(clientID uuid.UUID) (client models.Client, err error) {
	return
}

func (r *Repository) CreateAccount(clientID uuid.UUID, account *models.Account) error {
	return nil
}

func (r *Repository) GetAccountByID(accountID uuid.UUID) (account models.Account, err error) {
	return
}

func (r *Repository) updateBalance(account *models.Account) error {
	return nil
}

func (r *Repository) CreateTransaction(transaction *models.Transaction) error {
	return nil
}

func (r *Repository) GetTransactionsByAccount(accountID uuid.UUID) (transactions []models.Transaction, err error) {
	return
}

func (r *Repository) Close() error {
	return nil
}
