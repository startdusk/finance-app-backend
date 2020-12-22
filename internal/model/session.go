package model

type DeviceID string

var NilDeviceID DeviceID

// Session is represent user's session
type Session struct {
	UserID       UserID   `db:"user_id"`
	DeviceID     DeviceID `db:"device_id"`
	RefreshToken string   `db:"refresh_token"`
	ExpiresAt    int64    `db:"expires_at"`
}

// SessionData used to represent data sent in json body with requests
type SessionData struct {
	DeviceID DeviceID `db:"deviceID"`
}
