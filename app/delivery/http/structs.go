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

func (oClient outClient) FromModel(client models.Client, accounts []models.Account) outClient {
	oClient.ID = client.GetID().String()
	oClient.Name = client.Name
	oClient.Accounts = make([]outAccount, 0)
	for _, acc := range accounts {
		outAcc := outAccount.FromModel(outAccount{}, acc)
		oClient.Accounts = append(oClient.Accounts, outAcc)
	}
	return oClient
}

type outAccount struct {
	ID       string          `json:"id"`
	Currency outCurrency     `json:"currency"`
	Balance  decimal.Decimal `json:"balance"`
}

func (oAcc outAccount) FromModel(acc models.Account) outAccount {
	oAcc.ID = acc.GetID().String()
	oAcc.Balance = acc.Balance
	oAcc.Currency = outCurrency.FromModel(outCurrency{}, acc.Currency)
	return oAcc
}

type outCurrency struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
}

func (oCurrency outCurrency) FromModel(currency models.Currency) outCurrency {
	oCurrency.ID = currency.ID
	oCurrency.Code = currency.Code
	return oCurrency
}

type newTransaction struct {
	TypeID   uint            `json:"type_id"`
	SourceID string          `json:"source_id,omitempty"`
	TargetID string          `json:"target_id,omitempty"`
	Amount   decimal.Decimal `json:"amount"`
}

func (nTA newTransaction) ToModel() models.Transaction {
	sourceID, _ := uuid.Parse(nTA.SourceID)
	targetID, _ := uuid.Parse(nTA.TargetID)
	return models.NewTransaction(nTA.TypeID, sourceID, targetID, nTA.Amount)
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
	oTA.ID = t.GetID().String()
	oTA.Type = t.Type.Name
	oTA.SourceID = t.Source.GetID().String()
	oTA.TargetID = t.Target.GetID().String()
	oTA.Amount = t.Amount
	return oTA
}
