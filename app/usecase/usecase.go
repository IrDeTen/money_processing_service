package usecase

import (
	"github.com/IrDeTen/money_processing_service.git/app"
	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
)

type Usecase struct {
	repo app.IRepository
}

func NewUsecase(repo app.IRepository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) CreateClient(name string) (clientID uuid.UUID, err error) {
	return uuid.UUID{}, nil
}

func (u *Usecase) GetClient(clientID uuid.UUID) (client *models.Client, err error) {
	return nil, nil
}

func (u *Usecase) CreateAccount(clientID uuid.UUID, currencyID uint) (accountID uuid.UUID, err error) {
	return uuid.UUID{}, nil

}

func (u *Usecase) GetAccount(accountID uuid.UUID) (account *models.Account, err error) {
	return nil, nil

}

func (u *Usecase) CreateTransaction(transaction models.Transaction) (transactionID uuid.UUID, err error) {
	return uuid.UUID{}, nil

}

func (u *Usecase) GetTransactions(accountID uuid.UUID) (list []models.Transaction, err error) {
	return nil, nil
}
