CREATE TABLE sessions (
	user_id UUID REFERENCES users,
	device_id TEXT,
	refresh_token TEXT,
	expires_at INTEGER,
	PRIMARY KEY (user_id, device_id)
);
