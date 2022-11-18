package models

type TransactionType struct {
	ID   uint
	Name string
}

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
