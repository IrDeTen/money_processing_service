package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	errInvalidAccountUUID = errors.New("invalid id for account")
)

type Account struct {
	id       uuid.UUID
	Currency Currency
	Balance  decimal.Decimal
	Client   Client
}

func CreateNewAccount(clientID uuid.UUID, currencyID uint) (Account, error) {
	var client Client
	if err := client.SetID(clientID); err != nil {
		return Account{}, err
	}
	account := Account{
		Currency: Currency{
			ID: currencyID,
		},
		Client: client,
	}
	return account, nil
}

func (a *Account) GetID() uuid.UUID {
	return a.id
}

func (a *Account) SetID(id uuid.UUID) error {
	if id == uuid.Nil {
		return errInvalidAccountUUID
	}
	a.id = id
	return nil
}

func (a *Account) Withdraw(amount decimal.Decimal) error {
	if a.Balance.Sub(amount).IsNegative() {
		return errInsufficientMoney
	}
	a.Balance = a.Balance.Sub(amount)
	return nil
}

func (a *Account) Deposit(amount decimal.Decimal) {
	a.Balance = a.Balance.Add(amount)
}
