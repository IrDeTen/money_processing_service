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
	errIvalidTransactionUUID     = errors.New("invalid id for transaction")
)

type Transaction struct {
	id           uuid.UUID
	Type         TransactionType
	Amount       decimal.Decimal
	Source       Account
	Target       Account
	CreationDate time.Time
}

func NewTransaction(typeID uint, sourceID, targetID uuid.UUID, amount decimal.Decimal) Transaction {
	var source, target Account
	if sourceID != uuid.Nil {
		source.SetID(sourceID)
	}
	if targetID != uuid.Nil {
		target.SetID(targetID)
	}
	return Transaction{
		Type: TransactionType{
			ID: typeID,
		},
		Amount:       amount,
		Source:       source,
		Target:       target,
		CreationDate: time.Now(),
	}
}

func (t *Transaction) GetID() uuid.UUID {
	return t.id
}

func (t *Transaction) SetID(id uuid.UUID) error {
	if id == uuid.Nil {
		return errIvalidTransactionUUID
	}
	t.id = id
	return nil
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
	t.Target.Deposit(t.Amount)
}

func (t *Transaction) Withdraw() error {
	return t.Source.Withdraw(t.Amount)
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
	return t.Type.ID == Withdraw.ID
}

func (t *Transaction) IsTransfer() bool {
	return t.Type.ID == Transfer.ID
}

func (t *Transaction) SetType(typeID uint) {
	t.Type = transactionTypeMap[typeID]
}
