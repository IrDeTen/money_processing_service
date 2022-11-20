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
	return u.repo.CreateClient(client)
}

func (u *Usecase) GetClient(clientID uuid.UUID) (client models.Client, accounts []models.Account, err error) {
	client, err = u.repo.GetClient(clientID)
	if err != nil {
		err = errClientCreationFailed
		return
	}
	accounts, err = u.repo.GetAccountsByClientID(clientID)
	if err != nil {
		return
	}
	return
}

func (u *Usecase) CreateAccount(clientID uuid.UUID, currencyID uint) (uuid.UUID, error) {
	account, err := models.CreateNewAccount(clientID, currencyID)
	if err != nil {
		return uuid.Nil, err
	}
	return u.repo.CreateAccount(account)
}

func (u *Usecase) GetAccount(accountID uuid.UUID) (account models.Account, err error) {
	return u.repo.GetAccountByID(accountID)
}

func (u *Usecase) CreateTransaction(transaction models.Transaction) (transactionID uuid.UUID, err error) {
	if err = transaction.GeneralChecksForTransactionType(); err != nil {
		return
	}
	switch {
	case transaction.IsDeposit():
		var account models.Account
		if account, err = u.deposit(&transaction); err != nil {
			return
		}
		return u.repo.CreateTransaction(transaction, account)
	case transaction.IsWithdraw():
		var account models.Account
		if account, err = u.withdraw(&transaction); err != nil {
			return
		}
		return u.repo.CreateTransaction(transaction, account)
	case transaction.IsTransfer():
		var accounts []models.Account
		if accounts, err = u.transfer(&transaction); err != nil {
			return
		}
		return u.repo.CreateTransaction(transaction, accounts...)
	}

	return uuid.UUID{}, errInvalidTransactionType
}

func (u *Usecase) deposit(t *models.Transaction) (account models.Account, err error) {
	if t.Target, err = u.repo.GetAccountByID(t.Target.GetID()); err != nil {
		return models.Account{}, errAccountDoesntExist
	}
	t.Deposit()
	return t.Target, nil
}

func (u *Usecase) withdraw(t *models.Transaction) (account models.Account, err error) {
	if t.Source, err = u.repo.GetAccountByID(t.Source.GetID()); err != nil {
		return models.Account{}, errAccountDoesntExist
	}
	if err = t.Withdraw(); err != nil {
		return
	}
	return t.Source, nil
}

func (u *Usecase) transfer(t *models.Transaction) (accounts []models.Account, err error) {
	if t.Target, err = u.repo.GetAccountByID(t.Target.GetID()); err != nil {
		return nil, errAccountDoesntExist
	}
	accounts = append(accounts, t.Target)
	if t.Source, err = u.repo.GetAccountByID(t.Source.GetID()); err != nil {
		return nil, errAccountDoesntExist
	}
	accounts = append(accounts, t.Source)
	if err = t.Transfer(); err != nil {
		return
	}
	return
}

func (u *Usecase) GetTransactions(accountID uuid.UUID) (list []models.Transaction, err error) {
	return u.repo.GetTransactionsByAccount(accountID)
}
