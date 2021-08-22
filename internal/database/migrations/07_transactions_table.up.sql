CREATE TYPE transaction_type AS ENUM (
	'income',
	'expense'
);

CREATE TABLE transactions (
	transaction_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL REFERENCES users,
	account_id UUID NOT NULL REFERENCES accounts,
	category_id UUID NOT NULL REFERENCES categories,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP,

	date TIMESTAMP NOT NULL,
	type transaction_type NOT NULL,
	amount INTEGER NOT NULL,
	notes TEXT NOT NULL DEFAULT ''
);
