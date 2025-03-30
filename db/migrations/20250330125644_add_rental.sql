-- migrate:up
CREATE TABLE rentals (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL REFERENCES users(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    rental_type VARCHAR(100) NOT NULL,
    camper_id INT NOT NULL REFERENCES campers(id),
    driver_id INT REFERENCES drivers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE rental_equipments (
    rental_id INT NOT NULL REFERENCES rentals(id) ON DELETE CASCADE,
    equipment_id INT NOT NULL REFERENCES equipments(id) ON DELETE CASCADE,
    PRIMARY KEY (rental_id, equipment_id)
);

-- migrate:down
DROP TABLE IF EXISTS rental_equipments;
DROP TABLE IF EXISTS rentals;
