package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	errZeroValueTransaction  = errors.New("zero-value transaction")
	errInsufficientMoney     = errors.New("insufficient money in the account")
	errDifferentCurrencies   = errors.New("accounts with different currencies")
	errIvalidTransactionUUID = errors.New("invalid id for transaction")
)

type Transaction struct {
	id           uuid.UUID
	Type         TransactionType
	Amount       decimal.Decimal
	Source       Account
	Target       Account
	CreationDate time.Time
}

func NewTransaction(typeID uint, sourceID, targetID uuid.UUID, amount decimal.Decimal) (t Transaction, err error) {
	var source, target Account
	if sourceID != uuid.Nil {
		source.SetID(sourceID)
	}
	if targetID != uuid.Nil {
		target.SetID(targetID)
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return t, errZeroValueTransaction
	}
	t = Transaction{
		Type: TransactionType{
			ID: typeID,
		},
		Amount:       amount,
		Source:       source,
		Target:       target,
		CreationDate: time.Now(),
	}
	return
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

// The method is used to calculate the new account balance when making a deposit
func (t *Transaction) Deposit() {
	t.Target.Deposit(t.Amount)
}

// The method is used to calculate the new account balance when withdrawing
func (t *Transaction) Withdraw() error {
	return t.Source.Withdraw(t.Amount)
}

// The method is used to calculate the new balance of accounts during the transfer
func (t *Transaction) Transfer() error {
	if t.Source.Currency != t.Target.Currency {
		return errDifferentCurrencies
	}
	t.Deposit()
	return t.Withdraw()
}

// Returns true if transaction type is deposit
func (t *Transaction) IsDeposit() bool {
	return t.Type.ID == Deposit.ID
}

// Returns true if transaction type is withdraw
func (t *Transaction) IsWithdraw() bool {
	return t.Type.ID == Withdraw.ID
}

// Returns true if transaction type is transfer
func (t *Transaction) IsTransfer() bool {
	return t.Type.ID == Transfer.ID
}

func (t *Transaction) SetType(typeID uint) {
	t.Type = transactionTypeMap[typeID]
}
