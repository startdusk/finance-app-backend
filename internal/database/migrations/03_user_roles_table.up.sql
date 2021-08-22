
-- Create type Roles to avoid incorect input
-- We need only 'admin' role
-- Role 'member' is if user exists in database
-- If we will need more roles we will add them to this ENUM
CREATE TYPE user_role AS ENUM (
	'admin'
);

CREATE TABLE user_roles (
	user_id UUID REFERENCES users,
	role user_role NOT NULL,
	PRIMARY KEY(user_id, role)
);

-- Create index for roles
CREATE INDEX user_roles_user
	ON user_roles (user_id)
