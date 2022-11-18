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

func (u *Usecase) CreateClient(client models.Client) (uuid.UUID, error) {
	return client.GetID(), u.repo.CreateClient(&client)
}

func (u *Usecase) GetClient(clientID uuid.UUID) (client models.Client, accounts []models.Account, err error) {
	client, err = u.repo.GetClient(clientID)
	return
}

func (u *Usecase) CreateAccount(clientID uuid.UUID, currencyID uint) (uuid.UUID, error) {
	account := models.CreateNewAccount(currencyID)
	err := u.repo.CreateAccount(clientID, &account)
	return account.GetID(), err
}

func (u *Usecase) GetAccount(accountID uuid.UUID) (account models.Account, err error) {
	return u.repo.GetAccountByID(accountID)
}

func (u *Usecase) CreateTransaction(transaction models.Transaction) (transactionID uuid.UUID, err error) {
	if err = transaction.GeneralChecksForTransactionType(); err != nil {
		return
	}
	if transaction.IsDeposit() {
		if err = u.deposit(transaction); err != nil {
			return
		}
	}
	if transaction.IsWithdraw() {
		if err = u.withdraw(transaction); err != nil {
			return
		}
	}
	if transaction.IsTransfer() {
		if err = u.transfer(transaction); err != nil {
			return
		}
	}
	return transaction.GetID(), u.repo.CreateTransaction(&transaction)
}

func (u *Usecase) deposit(t models.Transaction) (err error) {
	if t.Target, err = u.repo.GetAccountByID(t.Target.GetID()); err != nil {
		return errAccountDoesntExist
	}
	t.Deposit()
	return nil
}

func (u *Usecase) withdraw(t models.Transaction) (err error) {
	if t.Source, err = u.repo.GetAccountByID(t.Source.GetID()); err != nil {
		return errAccountDoesntExist
	}
	if err = t.Withdraw(); err != nil {
		return
	}
	return nil
}

func (u *Usecase) transfer(t models.Transaction) (err error) {
	if t.Target, err = u.repo.GetAccountByID(t.Target.GetID()); err != nil {
		return errAccountDoesntExist
	}
	if t.Source, err = u.repo.GetAccountByID(t.Source.GetID()); err != nil {
		return errAccountDoesntExist
	}
	if err = t.Transfer(); err != nil {
		return
	}
	return nil
}

func (u *Usecase) GetTransactions(accountID uuid.UUID) (list []models.Transaction, err error) {
	return u.repo.GetTransactionsByAccount(accountID)
}
