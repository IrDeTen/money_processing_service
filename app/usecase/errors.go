package usecase

import "errors"

var (
	errClientCreationFailed     = errors.New("client creation failed")
	errAccountDoesntExist     = errors.New("specified account doesn't exist")
	errInvalidTransactionType = errors.New("invalid transaction type")
)
