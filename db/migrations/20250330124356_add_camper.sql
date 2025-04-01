-- migrate:up
CREATE TABLE campers (
    id VARCHAR(255) PRIMARY KEY,
    image_url TEXT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    license_plate VARCHAR(50) UNIQUE NOT NULL,
    year INT NOT NULL,
    capacity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    condition VARCHAR(50) NOT NULL,
    last_maintenance TIMESTAMP,
    transmission VARCHAR(50),
    fuel_type VARCHAR(50),
    drivetrain_config VARCHAR(50),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS campers;