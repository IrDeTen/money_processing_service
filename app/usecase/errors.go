package usecase

import "errors"

var (
	errAccountDoesntExist     = errors.New("specified account doesn't exist")
	errInvalidTransactionType = errors.New("invalid transaction type")
)
