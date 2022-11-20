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

func (u *Usecase) CreateClient(client models.Client) (id uuid.UUID, err error) {
	return u.repo.CreateClient(client)
}

func (u *Usecase) GetClient(clientID uuid.UUID) (client models.Client, accounts []models.Account, err error) {

	if client, err = u.repo.GetClient(clientID); err != nil {
		return
	}
	if accounts, err = u.repo.GetAccountsByClientID(clientID); err != nil {
		return
	}
	return
}

func (u *Usecase) CreateAccount(clientID uuid.UUID, currencyID uint) (id uuid.UUID, err error) {

	if _, err = u.repo.GetClient(clientID); err != nil {
		err = errClientDoesntExist
		return
	}

	if _, err = u.repo.GetCurrencyByID(currencyID); err != nil {
		err = errCurrencyDoesntExist
		return
	}

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
	var (
		account  models.Account
		accounts []models.Account
	)

	switch {
	case transaction.IsDeposit():
		if account, err = u.deposit(&transaction); err != nil {
			return
		}
		accounts = append(accounts, account)

	case transaction.IsWithdraw():
		if account, err = u.withdraw(&transaction); err != nil {
			return
		}
		accounts = append(accounts, account)

	case transaction.IsTransfer():
		if accounts, err = u.transfer(&transaction); err != nil {
			return
		}

	default:
		return uuid.Nil, errInvalidTransactionType

	}

	return u.repo.CreateTransaction(transaction, accounts...)
}

// The function checks the account ID, obtains account information and calculates a new balance.
// Returns an account object to update its balance in the database and an error if it occurred
func (u *Usecase) deposit(t *models.Transaction) (account models.Account, err error) {
	if t.Target.GetID() == uuid.Nil || t.Source.GetID() != uuid.Nil{
		return models.Account{}, errInvalidAccountForDeposit
	}
	if t.Target, err = u.repo.GetAccountByID(t.Target.GetID()); err != nil {
		return models.Account{}, errAccountDoesntExist
	}
	t.Deposit()
	return t.Target, nil
}

// The function checks the account ID, obtains account information and calculates a new balance.
// Returns an account object to update its balance in the database and an error if it occurred
func (u *Usecase) withdraw(t *models.Transaction) (account models.Account, err error) {
	if t.Source.GetID() == uuid.Nil || t.Target.GetID() != uuid.Nil{
		return models.Account{}, errInvalidAccountForWithdraw
	}

	if t.Source, err = u.repo.GetAccountByID(t.Source.GetID()); err != nil {
		return models.Account{}, errAccountDoesntExist
	}
	if err = t.Withdraw(); err != nil {
		return
	}
	return t.Source, nil
}

// The function checks the account ID, the correspondence of currencies, obtains account information and calculates a new balance.
// Returns both accounts to update their balances in the database and an error if it occurred
func (u *Usecase) transfer(t *models.Transaction) (accounts []models.Account, err error) {
	if t.Source.GetID() == uuid.Nil || t.Target.GetID() == uuid.Nil {
		return nil, errInvalidAccountForTransfer
	}
	if t.Source.GetID() == t.Target.GetID() {
		return nil, errSameAccount
	}
	if t.Target, err = u.repo.GetAccountByID(t.Target.GetID()); err != nil {
		return nil, errAccountDoesntExist
	}
	if t.Source, err = u.repo.GetAccountByID(t.Source.GetID()); err != nil {
		return nil, errAccountDoesntExist
	}
	if err = t.Transfer(); err != nil {
		return
	}
	accounts = append(accounts, t.Target, t.Source)
	return
}

func (u *Usecase) GetTransactions(accountID uuid.UUID) (list []models.Transaction, err error) {
	if _, err = u.repo.GetAccountByID(accountID); err != nil {
		return nil, errAccountDoesntExist
	}
	return u.repo.GetTransactionsByAccount(accountID)
}
