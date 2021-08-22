CREATE TABLE merchants (
	merchant_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL REFERENCES users,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP,
	name TEXT NOT NULL DEFAULT ''
);
