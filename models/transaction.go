package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	errInsufficientMoney         = errors.New("insufficient money in the account")
	errInvalidTransactionType    = errors.New("invalid transaction type")
	errInvalidAccountForDeposit  = errors.New("invalid account for deposit")
	errInvalidAccountForWithdraw = errors.New("invalid account for withdraw")
	errInvalidAccountForTransfer = errors.New("invalid account for transfer")
	errSameAccount               = errors.New("the same account for the transfer is specified")
	errDifferentCurrencies       = errors.New("accounts with different currencies")
)

type Transaction struct {
	id           uuid.UUID
	Type         TransactionType
	Amount       decimal.Decimal
	Source       Account
	Target       Account
	CreationDate time.Time
}

func NewTransaction() Transaction {
	return Transaction{
		id:           uuid.New(),
		CreationDate: time.Now(),
	}
}

func (t *Transaction) GetID() uuid.UUID {
	return t.id
}

func (t *Transaction) GeneralChecksForTransactionType() error {
	switch t.Type.ID {
	case Deposit.ID:
		if t.Target.GetID() == uuid.Nil || t.Source.GetID() != uuid.Nil {
			return errInvalidAccountForDeposit
		}

	case Withdraw.ID:
		if t.Source.GetID() == uuid.Nil || t.Target.GetID() != uuid.Nil {
			return errInvalidAccountForWithdraw
		}

	case Transfer.ID:
		if t.Source.GetID() == uuid.Nil || t.Target.GetID() == uuid.Nil {
			return errInvalidAccountForTransfer
		}
		if t.Source.GetID() == t.Target.GetID() {
			return errSameAccount
		}
	default:
		return errInvalidTransactionType
	}
	return nil
}

func (t *Transaction) Deposit() {
	t.Target.Balance.Add(t.Amount)
}

func (t *Transaction) Withdraw() error {
	newBalance := t.Target.Balance.Sub(t.Amount)
	if newBalance.IsNegative() {
		return errInsufficientMoney
	}
	t.Source.Balance = newBalance
	return nil
}

func (t *Transaction) Transfer() error {
	if t.Source.Currency != t.Target.Currency {
		return errDifferentCurrencies
	}
	t.Deposit()
	return t.Withdraw()
}

func (t *Transaction) IsDeposit() bool {
	return t.Type.ID == Deposit.ID
}

func (t *Transaction) IsWithdraw() bool {
	return t.Type.ID == Deposit.ID
}

func (t *Transaction) IsTransfer() bool {
	return t.Type.ID == Deposit.ID
}
