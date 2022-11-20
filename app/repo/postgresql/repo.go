package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/IrDeTen/money_processing_service.git/models"
	"github.com/IrDeTen/money_processing_service.git/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db        *sqlx.DB
	converter converter
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db:        sqlx.NewDb(db, "postgres"),
		converter: converter{},
	}
}

func (r *Repository) CreateClient(mClient models.Client) (id uuid.UUID, err error) {
	client := r.converter.ClientFromModel(mClient)
	query := "INSERT INTO clients(name) Values($1) RETURNING id"

	err = r.db.Get(&id, query, client.Name)
	if err != nil {
		logger.LogError(
			"create client",
			"app/repo/postgres/repo",
			"CreateClient",
			fmt.Sprintf("name: %s", client.Name),
			err)
		return
	}
	return
}

func (r *Repository) GetClient(clientID uuid.UUID) (client models.Client, err error) {
	var dbClient dbClient
	query := "SELECT * FROM clients WHERE id=$1"
	if err = r.db.Get(&dbClient, query, clientID); err != nil {
		logger.LogError(
			"get client",
			"app/repo/postgres/repo",
			"GetClient",
			fmt.Sprintf("client id: %s", clientID),
			err)
		return
	}
	client = r.converter.ClientToModel(dbClient)
	return
}

func (r *Repository) CreateAccount(account models.Account) (id uuid.UUID, err error) {
	dbAccount, err := r.converter.AccountFromModel(account)
	if err != nil {
		logger.LogError(
			"convert account to DB struct",
			"app/repo/postgres/repo",
			"CreateAccount",
			fmt.Sprintf("client id: %s, currency id: %d", dbAccount.ClientID.String(), dbAccount.CurrencyID),
			err)
		return
	}
	query := `INSERT INTO accounts(client_id, currency_id, balance) 
					Values(:client_id, :currency_id, :balance) 
				RETURNING id`
	prep, err := r.db.PrepareNamed(query)
	if err != nil {
		logger.LogError(
			"prepeare query for create account",
			"app/repo/postgres/repo",
			"CreateAccount",
			fmt.Sprintf("client id: %s, currency id: %d", dbAccount.ClientID.String(), dbAccount.CurrencyID),
			err)
		return
	}

	err = prep.Get(&id, dbAccount)
	if err != nil {
		logger.LogError(
			"create account",
			"app/repo/postgres/repo",
			"CreateAccount",
			fmt.Sprintf("client id: %s, currency id: %d", dbAccount.ClientID.String(), dbAccount.CurrencyID),
			err)
		return
	}
	return
}

func (r *Repository) GetAccountByID(accountID uuid.UUID) (account models.Account, err error) {
	var dbAccount dbAccount
	query := `SELECT A.*, C.code AS currency_code FROM accounts AS A 
				JOIN currencies AS C ON C.id = A.currency_id
			WHERE A.id =$1`
	if err = r.db.Get(&dbAccount, query, accountID); err != nil {
		logger.LogError(
			"get account",
			"app/repo/postgres/repo",
			"GetAccountByID",
			fmt.Sprintf("account id: %s", accountID.String()),
			err)
		return
	}
	account = r.converter.AccountToModel(dbAccount)
	return
}

func (r *Repository) GetAccountsByClientID(clientID uuid.UUID) (accounts []models.Account, err error) {
	dbAccounts := make([]dbAccount, 0)
	query := `SELECT A.*, C.code AS currency_code FROM accounts AS A 
				JOIN currencies AS C ON C.id = A.currency_id
			WHERE A.client_id =$1`
	if err = r.db.Select(&dbAccounts, query, clientID); err != nil {
		logger.LogError(
			"get accounts by client's id",
			"app/repo/postgres/repo",
			"GetAccountsByClientID",
			fmt.Sprintf("client id: %s", clientID),
			err)
		return
	}

	for _, val := range dbAccounts {
		acc := r.converter.AccountToModel(val)
		accounts = append(accounts, acc)
	}
	return
}

func (r *Repository) GetCurrencyByID(currencyID uint) (currency models.Currency, err error) {
	var curr dbCurrency
	query := `SELECT * FROM currencies WHERE id =$1`
	if err = r.db.Get(&curr, query, currencyID); err != nil {
		logger.LogError(
			"get currency",
			"app/repo/postgres/repo",
			"GetCurrencyByID",
			fmt.Sprintf("currency id: %d", currencyID),
			err)
		return
	}
	currency = r.converter.CurrencyToModel(curr)
	return

}

// Updating the balances of the specified accounts and creating a record of the transaction in the database
func (r *Repository) CreateTransaction(transaction models.Transaction, accounts ...models.Account) (id uuid.UUID, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		logger.LogError(
			"start transaction",
			"app/repo/postgres/repo",
			"CreateTransaction",
			fmt.Sprintf("transaction type id: %d, source account id: %s, target account id: %s",
				transaction.Type.ID, transaction.Source.GetID().String(), transaction.Target.GetID().String()),
			err)
		return
	}

	for _, account := range accounts {
		if err = r.updateBalance(tx, account); err != nil {
			return
		}
	}

	if id, err = r.writeTransaction(tx, transaction); err != nil {
		return
	}
	tx.Commit()
	return
}

// Update balance of the specified account in the database. Rolls back the transaction if errors occur in the process.
func (r *Repository) updateBalance(tx *sqlx.Tx, account models.Account) (err error) {
	var dbAccount dbAccount
	dbAccount, err = r.converter.AccountFromModel(account)
	if err != nil {
		logger.LogError(
			"convert account to DB structure",
			"app/repo/postgres/repo",
			"updateBalance",
			fmt.Sprintf("account id: %s", dbAccount.ID.String()),
			err)
		tx.Rollback()
		return
	}

	query := "UPDATE accounts SET balance = :balance WHERE id = :id"
	_, err = tx.NamedExec(query, dbAccount)
	if err != nil {
		logger.LogError(
			"account balance update",
			"app/repo/postgres/repo",
			"updateBalance",
			fmt.Sprintf("account id: %s", dbAccount.ID.String()),
			err)
		tx.Rollback()
		return
	}
	return
}

// Records a transaction between clients in the database. Rolls back the database transaction if errors occur in the process.
func (r *Repository) writeTransaction(tx *sqlx.Tx, transaction models.Transaction) (id uuid.UUID, err error) {
	dbTA := r.converter.TransactionFromModel(transaction)
	query := `INSERT INTO transactions(type_id, source_id, target_id, amount, creation_date) 
				VALUES(:type_id, :source_id, :target_id, :amount, :creation_date)
			RETURNING id`
	prep, err := r.db.PrepareNamed(query)
	if err != nil {
		logger.LogError(
			"prepeare query for create transaction",
			"app/repo/postgres/repo",
			"writeTransaction",
			fmt.Sprintf("transaction type id: %d, source account id: %s, target account id: %s",
				dbTA.TypeID, dbTA.SourceID.String, dbTA.TargetID.String),
			err)
		tx.Rollback()
		return
	}

	if err = prep.Get(&id, dbTA); err != nil {
		logger.LogError(
			"create transaction",
			"app/repo/postgres/repo",
			"writeTransaction",
			fmt.Sprintf("transaction type id: %d, source account id: %s, target account id: %s",
				dbTA.TypeID, dbTA.SourceID.String, dbTA.TargetID.String),
			err)
		tx.Rollback()
		return
	}
	return
}

func (r *Repository) GetTransactionsByAccount(accountID uuid.UUID) (transactions []models.Transaction, err error) {
	dbTAList := make([]dbTransaction, 0)
	query := "SELECT * FROM transactions WHERE source_id = $1 OR target_id = $1"
	err = r.db.Select(&dbTAList, query, accountID)
	if err != nil {
		logger.LogError(
			"get transactions by account id",
			"app/repo/postgres/repo",
			"GetTransactionsByAccount",
			fmt.Sprintf("account id: %s", accountID.String()),
			err)
		return
	}
	for _, val := range dbTAList {
		mTA := r.converter.TransactionToModel(val)
		transactions = append(transactions, mTA)
	}
	return
}

func (r *Repository) Close() error {
	return r.db.Close()
}
