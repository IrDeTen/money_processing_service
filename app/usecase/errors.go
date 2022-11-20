package usecase

import "errors"

var (
	errClientDoesntExist = errors.New("the client with the specified ID does not exist")

	errAccountDoesntExist = errors.New("specified account doesn't exist")

	errCurrencyDoesntExist = errors.New("invalid currency")

	errInvalidTransactionType    = errors.New("invalid transaction type")
	errInvalidAccountForDeposit  = errors.New("invalid account for deposit")
	errInvalidAccountForWithdraw = errors.New("invalid account for withdraw")
	errInvalidAccountForTransfer = errors.New("invalid account for transfer")
	errSameAccount               = errors.New("the same account for the transfer is specified")
)
