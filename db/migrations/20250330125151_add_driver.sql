-- migrate:up
CREATE TABLE drivers (
    id SERIAL PRIMARY KEY,
    photo TEXT,
    name VARCHAR(255) NOT NULL,
    license_number VARCHAR(100) UNIQUE NOT NULL,
    license_expiry DATE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL,
    join_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS drivers;
