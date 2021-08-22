package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/startdusk/finance-app-backend/internal/model"
)

// UsersDB persist users
type UsersDB interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID model.UserID) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	ListUsers(ctx context.Context) ([]*model.User, error)
	DeleteUser(ctx context.Context, userID model.UserID) (bool, error)
}

var ErrUserExist = errors.New("user with that email exists")

const createUserQuery = `
	INSERT INTO users (
		email, password_hash
	)
	VALUES (
		:email, :password_hash
	)
	RETURNING user_id;
`

func (d *database) CreateUser(ctx context.Context, user *model.User) error {
	rows, err := d.conn.NamedQueryContext(ctx, createUserQuery, user)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == UniqueViolation {
				if pqError.Constraint == "user_email" {
					return ErrUserExist
				}
			}
		}
		return errors.Wrap(err, "could not create user")
	}

	rows.Next()
	if err := rows.Scan(&user.ID); err != nil {
		return errors.Wrap(err, "could not get created user id")
	}
	return nil
}

const getUserByIDQuery = `
	SELECT user_id, email, password_hash, created_at 
	FROM users 
	WHERE user_id = $1 AND deleted_at IS NULL;
`

func (d *database) GetUserByID(ctx context.Context, userID model.UserID) (*model.User, error) {
	var user model.User
	if err := d.conn.GetContext(ctx, &user, getUserByIDQuery, userID); err != nil {
		return nil, err
	}
	return &user, nil
}

const getUserByEmailQuery = `
	SELECT user_id, email, password_hash, created_at 
	FROM users 
	WHERE email = $1 AND deleted_at IS NULL;
`

func (d *database) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := d.conn.GetContext(ctx, &user, getUserByEmailQuery, email); err != nil {
		return nil, err
	}
	return &user, nil
}

const listUserQuery = `
	SELECT user_id, email, password_hash, created_at 
	FROM users    
	WHERE deleted_at IS NULL;
`

func (d *database) ListUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := d.conn.SelectContext(ctx, &users, listUserQuery); err != nil {
		return nil, errors.Wrap(err, "could not get users")
	}

	return users, nil
}

// we don't delete records from database we want them as deleted by setting deleted_at time
// I do this to avoid errors by deleting records
// I change users email as well because I allow only unique emails
// Maybe in future we will lock account and not delete it.
// If we lock account we can save all users data and if user delete his account
// and regret after sometime we can unlock it
// TODO: check comment
const deleteUserQuery = `
	UPDATE users  
	SET deleted_at = NOW(),
		email = CONCAT(email, '-DELETE-', uuid_generate_v4()) 
	WHERE user_id = $1;
`

func (d *database) DeleteUser(ctx context.Context, userID model.UserID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}

	return true, nil
}

const updateUserQuery = `
	UPDATE users 
	SET password_hash = :password_hash 
	WHERE user_id = :user_id;
`

func (d *database) UpdateUser(ctx context.Context, user *model.User) error {
	result, err := d.conn.NamedExecContext(ctx, updateUserQuery, user)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("user not found")
	}

	return nil
}
