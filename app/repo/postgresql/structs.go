package postgresql

import (
	"database/sql"
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
	ID           uuid.UUID       `db:"id"`
	ClientID     uuid.UUID       `db:"client_id"`
	CurrencyID   uint            `db:"currency_id"`
	CurrencyCode string          `db:"currency_code"`
	Balance      decimal.Decimal `db:"balance"`
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
	mAcc.Currency.Code = dbAcc.CurrencyCode
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
	ID           uuid.UUID       `db:"id"`
	TypeID       uint            `db:"type_id"`
	Amount       decimal.Decimal `db:"amount"`
	SourceID     sql.NullString  `db:"source_id"`
	TargetID     sql.NullString  `db:"target_id"`
	CreationDate time.Time       `db:"creation_date"`
}

func (c converter) TransactionFromModel(mTA models.Transaction) (dbTA dbTransaction) {
	if mTA.GetID() != uuid.Nil {
		dbTA.ID = mTA.GetID()
	}
	if mTA.Source.GetID() != uuid.Nil {
		dbTA.SourceID = sql.NullString{
			Valid:  true,
			String: mTA.Source.GetID().String(),
		}
	}
	if mTA.Target.GetID() != uuid.Nil {
		dbTA.TargetID = sql.NullString{
			Valid:  true,
			String: mTA.Source.GetID().String(),
		}
	}
	dbTA.CreationDate = mTA.CreationDate
	dbTA.Amount = mTA.Amount
	dbTA.TypeID = mTA.Type.ID
	return
}

func (c converter) TransactionToModel(dbTA dbTransaction) (mTA models.Transaction) {
	mTA.SetID(dbTA.ID)
	if dbTA.SourceID.Valid {
		sourceID, _ := uuid.Parse(dbTA.SourceID.String)
		mTA.Source.SetID(sourceID)
	}
	if dbTA.TargetID.Valid {
		targetID, _ := uuid.Parse(dbTA.TargetID.String)
		mTA.Target.SetID(targetID)
	}
	mTA.SetType(dbTA.TypeID)
	mTA.Amount = dbTA.Amount
	mTA.CreationDate = dbTA.CreationDate
	return
}
