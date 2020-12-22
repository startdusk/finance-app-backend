package database

import (
	"context"

	"github.com/startdusk/finance-app-backend/internal/model"
)

type SessionDB interface {
	SaveRefreshToken(ctx context.Context, session model.Session) error
	GetSession(ctx context.Context, session model.Session) (*model.Session, error)
}

var insertOrUpdateSession = `
	INSERT INTO sessions (user_id, device_id, refresh_token, expires_at) 
	VALUES (:user_id, :device_id, :refresh_token, :expires_at) 

	ON CONFLICT (user_id, device_id) 
	DO
		UPDATE 
			SET refresh_token = :refresh_token,
				expires_at = :expires_at
`

func (d *database) SaveRefreshToken(ctx context.Context, session model.Session) error {
	if _, err := d.conn.NamedQueryContext(ctx, insertOrUpdateSession, session); err != nil {
		return err
	}
	return nil
}

var getSessionQuery = `
	SELECT user_id, device_id, refresh_token, expires_at 
	FROM sessions 
	WHERE user_id = $1 
		AND device_id = $2 
		AND refresh_token = $3 
		AND to_timestamp(expires_at) > NOW() 
`

func (d *database) GetSession(ctx context.Context, data model.Session) (*model.Session, error) {
	var session model.Session
	if err := d.conn.GetContext(ctx, &session, getSessionQuery, data.UserID, data.DeviceID, data.RefreshToken); err != nil {
		return nil, err
	}

	return &session, nil
}
