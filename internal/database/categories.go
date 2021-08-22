package database

import (
	"context"

	"github.com/pkg/errors"

	"github.com/startdusk/finance-app-backend/internal/model"
)

type CategoryDB interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, categoryID model.CategoryID) (*model.Category, error)
	ListCategoriesByUserID(ctx context.Context, userID model.UserID) ([]*model.Category, error)
	DeleteCategory(ctx context.Context, categoryID model.CategoryID) (bool, error)
}

const createCategoryQuery = `
	INSERT INTO categories (parent_id, user_id, name) 
		VALUES (:parent_id, :user_id, :name) 
	RETURNING category_id;
`

func (d *database) CreateCategory(ctx context.Context, category *model.Category) error {
	rows, err := d.conn.NamedQueryContext(ctx, createCategoryQuery, category)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&category.ID); err != nil {
		return err
	}

	return nil
}

const updateCategoryQuery = `
	UPDATE categories 
	SET parent_id = :parent_id, 
		name = :name 
	WHERE category_id = :category_id;
`

func (d *database) UpdateCategory(ctx context.Context, category *model.Category) error {
	result, err := d.conn.NamedExecContext(ctx, updateCategoryQuery, category)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("category not found")
	}

	return nil
}

const getCategoryByIDQuery = `
	SELECT category_id, parent_id, user_id, name, created_at, deleted_at 
	FROM categories  
	WHERE category_id = $1 AND deleted_at IS NULL;
`

func (d *database) GetCategoryByID(ctx context.Context, categoryID model.CategoryID) (*model.Category, error) {
	var category model.Category
	if err := d.conn.GetContext(ctx, &category, getCategoryByIDQuery, categoryID); err != nil {
		return nil, errors.Wrap(err, "could not get category")
	}

	return &category, nil
}

const listCategoryByUserIDQuery = `
	SELECT category_id, parent_id, user_id, name, created_at, deleted_at 
	FROM categories  
	WHERE user_id = $1 AND deleted_at IS NULL;
`

func (d *database) ListCategoriesByUserID(ctx context.Context, userID model.UserID) ([]*model.Category, error) {
	var categories []*model.Category
	if err := d.conn.SelectContext(ctx, &categories, listCategoryByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "could not get user's categories")
	}

	return categories, nil
}

// we don't delete records from database we want them as deleted by setting deleted_at time
const deleteCategoryQuery = `
	UPDATE categories  
	SET deleted_at = NOW() 
	WHERE category_id = $1;
`

func (d *database) DeleteCategory(ctx context.Context, categoryID model.CategoryID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteCategoryQuery, categoryID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}
