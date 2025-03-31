-- migrate:up
CREATE TABLE campers (
    id VARCHAR(255) PRIMARY KEY,
    image_url TEXT,
    name VARCHAR(255) NOT NULL,
    license_plate VARCHAR(50) UNIQUE NOT NULL,
    year INT NOT NULL,
    capacity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    condition VARCHAR(50) NOT NULL,
    last_maintenance TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS campers;