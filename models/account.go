package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	id       uuid.UUID
	Currency Currency
	Balance  decimal.Decimal
}

func CreateNewAccount(currencyID uint) Account {
	return Account{
		id: uuid.New(),
		Currency: Currency{
			ID: currencyID,
		},
	}
}

func (a *Account) GetID() uuid.UUID {
	return a.id
}

func (a *Account) Withdraw(currencyID uint, amount decimal.Decimal) error {
	if a.Balance.Sub(amount).IsNegative() {
		return errInsufficientMoney
	}
	a.Balance.Sub(amount)
	return nil
}

func (a *Account) Deposit(currencyID uint, amount decimal.Decimal) {
	a.Balance.Add(amount)
}
