package models

import (
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type Account struct {
	ID       uuid.UUID
	Currency Currency
	Balance  decimal.Decimal
}
