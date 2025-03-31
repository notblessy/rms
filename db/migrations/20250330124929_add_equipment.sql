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

CREATE TABLE camper_equipments (
    rental_id INT NOT NULL REFERENCES rentals(id) ON DELETE CASCADE,
    equipment_id INT NOT NULL REFERENCES equipments(id) ON DELETE CASCADE,
    PRIMARY KEY (rental_id, equipment_id)
);

-- migrate:down
DROP TABLE IF EXISTS camper_equipments;
DROP TABLE IF EXISTS equipments;
