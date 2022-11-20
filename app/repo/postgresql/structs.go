package postgresql

import (
	"errors"
	"time"

	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type converter struct {
}

type dbClient struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (c converter) ClientFromModel(mCl models.Client) (dbCl dbClient) {
	dbCl.Name = mCl.Name
	if mCl.GetID() != uuid.Nil {
		dbCl.ID = mCl.GetID()
	}
	return
}

func (c converter) ClientToModel(dbClient dbClient) (mClient models.Client) {
	mClient.Name = dbClient.Name
	mClient.SetID(dbClient.ID)
	return
}

type dbAccount struct {
	ID         uuid.UUID       `db:"id"`
	ClientID   uuid.UUID       `db:"client_id"`
	CurrencyID uint            `db:"currency_id"`
	Balance    decimal.Decimal `db:"balance"`
}

func (c converter) AccountFromModel(mAcc models.Account) (dbAcc dbAccount, err error) {
	dbAcc = dbAccount{
		CurrencyID: mAcc.Currency.ID,
		Balance:    mAcc.Balance,
	}
	if mAcc.GetID() != uuid.Nil {
		dbAcc.ID = mAcc.GetID()
	}
	if mAcc.Client.GetID() == uuid.Nil {
		err = errors.New("account contains invalid client ID")
		return
	}
	dbAcc.ClientID = mAcc.Client.GetID()
	return
}

func (c converter) AccountToModel(dbAcc dbAccount) (mAcc models.Account) {
	if dbAcc.ID != uuid.Nil {
		mAcc.SetID(dbAcc.ID)
	}
	if dbAcc.ClientID != uuid.Nil {
		mAcc.Client.SetID(dbAcc.ClientID)
	}
	mAcc.Currency.ID = dbAcc.CurrencyID
	mAcc.Balance = dbAcc.Balance
	return
}

type dbCurrency struct {
	ID   uint   `db:"id"`
	Code string `db:"code"`
}

func (c converter) CurrencyFromModel(mCurr models.Currency) (dbCurr dbCurrency) {
	dbCurr.ID = mCurr.ID
	dbCurr.Code = mCurr.Code
	return
}

func (c converter) CurrencyToModel(dbCurr dbCurrency) (mCurr models.Currency) {
	mCurr.Code = dbCurr.Code
	mCurr.ID = dbCurr.ID
	return
}

type dbTransaction struct {
	ID           uuid.UUID       `json:"id"`
	TypeID       uint            `json:"type_id"`
	Amount       decimal.Decimal `json:"amount"`
	SourceID     uuid.UUID       `json:"source_id"`
	TargetID     uuid.UUID       `json:"target_id"`
	CreationDate time.Time       `json:"creation_date"`
}

func (c converter) TransactionFromModel(mTA models.Transaction) (dbTA dbTransaction) {
	if mTA.GetID() != uuid.Nil {
		dbTA.ID = mTA.GetID()
	}
	if mTA.Source.GetID() != uuid.Nil {
		dbTA.SourceID = mTA.Source.GetID()
	}
	if mTA.Target.GetID() != uuid.Nil {
		dbTA.TargetID = mTA.Target.GetID()
	}
	dbTA.CreationDate = mTA.CreationDate
	dbTA.Amount = mTA.Amount
	dbTA.TypeID = mTA.Type.ID
	return
}

func (c converter) TransactionToModel(dbTA dbTransaction) (mTA models.Transaction) {
	mTA.SetID(dbTA.ID)
	mTA.Source.SetID(dbTA.SourceID)
	mTA.Target.SetID(dbTA.TargetID)
	mTA.Type.ID = dbTA.TypeID
	mTA.Amount = dbTA.Amount
	mTA.CreationDate = dbTA.CreationDate
	return
}
