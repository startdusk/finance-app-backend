package database

import (
	"context"

	"github.com/pkg/errors"

	"github.com/startdusk/finance-app-backend/internal/model"
)

type MerchantDB interface {
	CreateMerchant(ctx context.Context, merchant *model.Merchant) error
	UpdateMerchant(ctx context.Context, merchant *model.Merchant) error
	GetMerchantByID(ctx context.Context, merchantID model.MerchantID) (*model.Merchant, error)
	ListMerchantsByUserID(ctx context.Context, userID model.UserID) ([]*model.Merchant, error)
	DeleteMerchant(ctx context.Context, merchantID model.MerchantID) (bool, error)
}

const createMerchantQuery = `
	INSERT INTO merchants (user_id, name) 
		VALUES (:user_id, :name) 
	RETURNING merchant_id;
`

func (d *database) CreateMerchant(ctx context.Context, merchant *model.Merchant) error {
	rows, err := d.conn.NamedQueryContext(ctx, createMerchantQuery, merchant)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&merchant.ID); err != nil {
		return err
	}

	return nil
}

const updateMerchantQuery = `
	UPDATE merchants 
	SET name = :name 
	WHERE merchant_id = :merchant_id;
`

func (d *database) UpdateMerchant(ctx context.Context, merchant *model.Merchant) error {
	result, err := d.conn.NamedExecContext(ctx, updateMerchantQuery, merchant)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("merchant not found")
	}

	return nil
}

const getMerchantByIDQuery = `
	SELECT merchant_id, user_id, name, created_at, deleted_at 
	FROM merchants   
	WHERE merchant_id = $1 AND deleted_at IS NULL;
`

func (d *database) GetMerchantByID(ctx context.Context, merchantID model.MerchantID) (*model.Merchant, error) {
	var merchant model.Merchant
	if err := d.conn.GetContext(ctx, &merchant, getMerchantByIDQuery, merchantID); err != nil {
		return nil, errors.Wrap(err, "could not get merchant")
	}

	return &merchant, nil
}

const listMerchantByUserIDQuery = `
	SELECT merchant_id, user_id, name, created_at, deleted_at 
	FROM merchants   
	WHERE user_id = $1 AND deleted_at IS NULL;
`

func (d *database) ListMerchantsByUserID(ctx context.Context, userID model.UserID) ([]*model.Merchant, error) {
	var merchants []*model.Merchant
	if err := d.conn.SelectContext(ctx, &merchants, listMerchantByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get user's merchants")
	}

	return merchants, nil
}

// we don't delete records from database we want them as deleted by setting deleted_at time
const deleteMerchantQuery = `
	UPDATE merchants  
	SET deleted_at = NOW() 
	WHERE merchant_id = $1;
`

func (d *database) DeleteMerchant(ctx context.Context, merchantID model.MerchantID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteMerchantQuery, merchantID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
