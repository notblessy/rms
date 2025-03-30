-- migrate:up
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(150) NOT NULL UNIQUE,
    name VARCHAR(150),
    picture TEXT,
    phone VARCHAR(20),
    address TEXT NOT NULL,
    id_number VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(150),
    picture TEXT,
    role VARCHAR(32) NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT users_email_idx UNIQUE (email)
);

-- migrate:down
DROP TABLE IF EXISTS users;
