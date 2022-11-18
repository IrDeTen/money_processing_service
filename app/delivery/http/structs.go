package http

import (
	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type newClient struct {
	Name string `json:"name"`
}

func (nClient newClient) ToModel() models.Client {
	return models.CreateNewClient(nClient.Name)
}

type outClient struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Accounts []outAccount `json:"accounts"`
}

func (oClient *outClient) FromModel(client models.Client, accounts []models.Account) {
	list := make([]outAccount, 0)
	for _, acc := range accounts {
		outAcc := outAccount.FromModel(outAccount{}, acc)
		list = append(list, outAcc)
	}
	oClient = &outClient{
		ID:       client.GetID().String(),
		Name:     client.Name,
		Accounts: list,
	}
}

type newAccount struct {
	ClientID   string `json:"client_id"`
	CurrencyID uint   `json:"currency_id"`
}

func (nAcc newAccount) ToModel() models.Account {
	return models.Account{}
}

type outAccount struct {
	ID       string          `json:"id"`
	Currency outCurrency     `json:"currency"`
	Balance  decimal.Decimal `json:"balance"`
}

func (oAcc outAccount) FromModel(acc models.Account) outAccount {
	return outAccount{
		ID:       acc.GetID().String(),
		Currency: outCurrency.FromModel(outCurrency{}, acc.Currency),
		Balance:  acc.Balance,
	}
}

type outCurrency struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
}

func (oCurrency outCurrency) FromModel(currency models.Currency) outCurrency {
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

func (nTA newTransaction) ToModel() (models.Transaction, error) {
	var sourceID, targetID uuid.UUID
	var err error

	sourceID, err = uuid.Parse(nTA.SourceID)
	if len(nTA.SourceID) > 0 && err != nil {
		return models.Transaction{}, errInvalidAccountID
	}

	targetID, err = uuid.Parse(nTA.TargetID)
	if len(nTA.TargetID) > 0 && err != nil {
		return models.Transaction{}, errInvalidAccountID
	}
	return models.NewTransaction(nTA.TypeID, sourceID, targetID, nTA.Amount), nil
}

type outTransaction struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	SourceID string          `json:"source_id,omitempty"`
	TargetID string          `json:"target_id,omitempty"`
	Amount   decimal.Decimal `json:"amount"`
	Date     string          `json:"date"`
}

func (oTA outTransaction) FromModel(t models.Transaction) outTransaction {
	return outTransaction{
		ID:       t.GetID().String(),
		Type:     t.Type.Name,
		SourceID: t.Source.GetID().String(),
		TargetID: t.Target.GetID().String(),
		Amount:   t.Amount,
		Date:     t.CreationDate.String(),
	}
}
