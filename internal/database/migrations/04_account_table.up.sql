CREATE TABLE accounts (
	account_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL REFERENCES users,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP,
	start_balance INTEGER NOT NULL,
	account_type TEXT NOT NULL,
	account_name TEXT NOT NULL,
	currency TEXT NOT NULL 
);
