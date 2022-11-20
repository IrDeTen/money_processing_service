package http

import (
	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type converter struct{}

type newClient struct {
	Name string `json:"name,omitempty"`
}

func (c converter) ClientToModel(client newClient) models.Client {
	return models.CreateNewClient(client.Name)
}

type outClient struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Accounts []outAccount `json:"accounts"`
}

func (c converter) ClientFromModel(client models.Client, accounts []models.Account) outClient {
	list := make([]outAccount, 0)
	for _, account := range accounts {
		list = append(list, c.AccountFromModel(account))
	}
	return outClient{
		ID:       client.GetID().String(),
		Name:     client.Name,
		Accounts: list,
	}
}

type newAccount struct {
	ClientID   string `json:"client_id"`
	CurrencyID uint   `json:"currency_id"`
}

type outAccount struct {
	ID       string          `json:"id"`
	Currency outCurrency     `json:"currency"`
	Balance  decimal.Decimal `json:"balance"`
}

func (c converter) AccountFromModel(account models.Account) outAccount {
	return outAccount{
		ID:       account.GetID().String(),
		Currency: c.CurrencyFromModel(account.Currency),
		Balance:  account.Balance,
	}
}

type outCurrency struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
}

func (c converter) CurrencyFromModel(currency models.Currency) outCurrency {
	return outCurrency{
		ID:   currency.ID,
		Code: currency.Code,
	}
}

type newTransaction struct {
	TypeID   uint            `json:"type_id"`
	SourceID string          `json:"source_id,omitempty"`
	TargetID string          `json:"target_id,omitempty"`
	Amount   decimal.Decimal `json:"amount"`
}

func (c converter) NewTransactionToModel(transaction newTransaction) (models.Transaction, error) {
	var sourceID, targetID uuid.UUID
	var err error
	sourceID, err = uuid.Parse(transaction.SourceID)
	if len(transaction.SourceID) > 0 && err != nil {
		return models.Transaction{}, errInvalidAccountID
	}

	targetID, err = uuid.Parse(transaction.TargetID)
	if len(transaction.TargetID) > 0 && err != nil {
		return models.Transaction{}, errInvalidAccountID
	}
	return models.NewTransaction(transaction.TypeID, sourceID, targetID, transaction.Amount)
}

type outTransaction struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	SourceID string          `json:"source_id,omitempty"`
	TargetID string          `json:"target_id,omitempty"`
	Amount   decimal.Decimal `json:"amount"`
	Date     string          `json:"date"`
}

func (c converter) TransactionFromModel(transaction models.Transaction) outTransaction {
	return outTransaction{
		ID:       transaction.GetID().String(),
		Type:     transaction.Type.Name,
		SourceID: transaction.Source.GetID().String(),
		TargetID: transaction.Target.GetID().String(),
		Amount:   transaction.Amount,
		Date:     transaction.CreationDate.String(),
	}
}
