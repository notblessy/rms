-- migrate:up
CREATE TABLE equipments (
    id SERIAL PRIMARY KEY,
    image_url TEXT,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    stock INT NOT NULL,
    condition VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS equipments;
