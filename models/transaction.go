package models

import (
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type Transaction struct {
	ID     uuid.UUID
	Type   TransactionType
	Amount decimal.Decimal
	Source Account
	Target Account
}

type TransactionType struct {
	ID   uint
	Name string
}

// TODO change work with transaction type
var (
	Deposit = TransactionType{
		ID:   1,
		Name: "Deposit",
	}

	Withdraw = TransactionType{
		ID:   2,
		Name: "Withdraw",
	}

	Transfer = TransactionType{
		ID:   3,
		Name: "Transfer",
	}
)
