package database

import (
	"context"

	"github.com/pkg/errors"

	"github.com/startdusk/finance-app-backend/internal/model"
)

type AccountDB interface {
	CreateAccount(ctx context.Context, account *model.Account) error
	UpdateAccount(ctx context.Context, account *model.Account) error
	GetAccountByID(ctx context.Context, accountID model.AccountID) (*model.Account, error)
	ListAccountsByUserID(ctx context.Context, userID model.UserID) ([]*model.Account, error)
	DeleteAccount(ctx context.Context, accountID model.AccountID) (bool, error)
}

const createAccountQuery = `
	INSERT INTO accounts (user_id, start_balance, account_type, account_name, currency) 
		VALUES (:user_id, :start_balance, :account_type, :account_name, :currency) 
	RETURNING account_id;
`

func (d *database) CreateAccount(ctx context.Context, account *model.Account) error {
	rows, err := d.conn.NamedQueryContext(ctx, createAccountQuery, account)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&account.ID); err != nil {
		return err
	}

	return nil
}

const updateAccountQuery = `
	UPDATE accounts 
	SET start_balance = :start_balance, 
		account_type = :account_type,
		account_name = :account_name,
		currency = :currency 
	WHERE account_id = :account_id;
`

func (d *database) UpdateAccount(ctx context.Context, account *model.Account) error {
	result, err := d.conn.NamedExecContext(ctx, updateAccountQuery, account)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("account not found")
	}

	return nil
}

const getAccountByIDQuery = `
	SELECT account_id, user_id, start_balance, account_type, account_name, currency, created_at, deleted_at 
	FROM accounts 
	WHERE account_id = $1;
`

func (d *database) GetAccountByID(ctx context.Context, accountID model.AccountID) (*model.Account, error) {
	var account model.Account
	if err := d.conn.GetContext(ctx, &account, getAccountByIDQuery, accountID); err != nil {
		return nil, errors.Wrap(err, "could not get account")
	}

	return &account, nil
}

const listAccountByUserIDQuery = `
	SELECT account_id, user_id, start_balance, account_type, account_name, currency, created_at, deleted_at 
	FROM accounts 
	WHERE user_id = $1 AND deleted_at IS NULL;
`

func (d *database) ListAccountsByUserID(ctx context.Context, userID model.UserID) ([]*model.Account, error) {
	var accounts []*model.Account
	if err := d.conn.SelectContext(ctx, &accounts, listAccountByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get user's accounts")
	}

	return accounts, nil
}

// we don't delete records from database we want them as deleted by setting deleted_at time
const deleteAccountQuery = `
	UPDATE accounts 
	SET deleted_at = NOW() 
	WHERE account_id = $1;
`

func (d *database) DeleteAccount(ctx context.Context, accountID model.AccountID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteAccountQuery, accountID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
