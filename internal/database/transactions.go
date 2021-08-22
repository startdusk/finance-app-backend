package database

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/startdusk/finance-app-backend/internal/model"
)

type TransactionDB interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetTransactionByID(ctx context.Context, transactionID model.TransactionID) (*model.Transaction, error)
	ListTransactionByCategoryID(ctx context.Context, categoryID model.CategoryID, from, to time.Time) ([]*model.Transaction, error)
	ListTransactionByAccountID(ctx context.Context, accountID model.AccountID, from, to time.Time) ([]*model.Transaction, error)
	ListTransactionByUserID(ctx context.Context, userID model.UserID, from, to time.Time) ([]*model.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID model.TransactionID) (bool, error)
}

const createTransactionQuery = `
	INSERT INTO transactions (user_id, account_id, category_id, date, type, amount, notes) 
		VALUES (:user_id, :account_id, :category_id, :date, :type, :amount, :notes) 
	RETURNING transaction_id;
`

func (d *database) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	rows, err := d.conn.NamedQueryContext(ctx, createTransactionQuery, transaction)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&transaction.ID); err != nil {
		return err
	}

	return nil
}

const updateTransactionQuery = `
	UPDATE transactions 
	SET account_id = :account_id, 
		category_id = :category_id, 
		date = :date, 
		type = :type, 
		amount = :amount, 
		notes = :notes 
	WHERE transaction_id = :transaction_id;
`

func (d *database) UpdateTransaction(ctx context.Context, transaction *model.Transaction) error {
	result, err := d.conn.NamedExecContext(ctx, updateTransactionQuery, transaction)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("transaction not found")
	}

	return nil
}

const getTransactionByIDQuery = `
	SELECT transaction_id, user_id, account_id, category_id, date, type, amount, notes, created_at, deleted_at 
	FROM transactions   
	WHERE transaction_id = $1 
		AND deleted_at IS NULL;
`

func (d *database) GetTransactionByID(ctx context.Context, transactionID model.TransactionID) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := d.conn.GetContext(ctx, &transaction, getTransactionByIDQuery, transactionID); err != nil {
		return nil, errors.Wrap(err, "could not get transaction")
	}

	return &transaction, nil
}

const listTransactionByUserIDQuery = `
	SELECT transaction_id, user_id, account_id, category_id, date, type, amount, notes, created_at, deleted_at 
	FROM transactions 
	WHERE user_id = $1 
		AND deleted_at IS NULL
		AND date > $2 
		AND date < $3;
`

func (d *database) ListTransactionByUserID(ctx context.Context, userID model.UserID, from, to time.Time) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	if err := d.conn.SelectContext(ctx, &transactions, listTransactionByUserIDQuery, userID, from, to); err != nil {
		return nil, errors.Wrap(err, "could not get user's transactions")
	}

	return transactions, nil
}

const listTransactionByCategoryIDQuery = `
	SELECT transaction_id, user_id, account_id, category_id, date, type, amount, notes, created_at, deleted_at 
	FROM transactions 
	WHERE category_id = $1 
		AND deleted_at IS NULL 
		AND date > $2 
		AND date < $3;
`

func (d *database) ListTransactionByCategoryID(ctx context.Context, categoryID model.CategoryID, from, to time.Time) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	if err := d.conn.SelectContext(ctx, &transactions, listTransactionByCategoryIDQuery, categoryID, from, to); err != nil {
		return nil, errors.Wrap(err, "could not get categories transactions")
	}

	return transactions, nil
}

const listTransactionByAccountIDQuery = `
	SELECT transaction_id, user_id, account_id, category_id, date, type, amount, notes, created_at, deleted_at 
	FROM transactions 
	WHERE account_id = $1 
		AND deleted_at IS NULL
		AND date > $2 
		AND date < $3;
`

func (d *database) ListTransactionByAccountID(ctx context.Context, accountID model.AccountID, from, to time.Time) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	if err := d.conn.SelectContext(ctx, &transactions, listTransactionByAccountIDQuery, accountID, from, to); err != nil {
		return nil, errors.Wrap(err, "could not get accounts transactions")
	}

	return transactions, nil
}

// we don't delete records from database we want them as deleted by setting deleted_at time
const deleteTransactionQuery = `
	UPDATE transactions  
	SET deleted_at = NOW() 
	WHERE transaction_id = $1;
`

func (d *database) DeleteTransaction(ctx context.Context, transactionID model.TransactionID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteTransactionQuery, transactionID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
