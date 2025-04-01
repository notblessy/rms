-- migrate:up
CREATE TABLE rentals (
    id VARCHAR(255) PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL REFERENCES users(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    rental_type VARCHAR(100) NOT NULL,
    camper_id VARCHAR(255) NOT NULL REFERENCES campers(id),
    driver_id VARCHAR(255) REFERENCES drivers(id),
    status VARCHAR(50) NOT NULL,
    grand_total DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) DEFAULT 0,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE rental_equipments (
    rental_id VARCHAR(255) NOT NULL REFERENCES rentals(id) ON DELETE CASCADE,
    equipment_id VARCHAR(255) NOT NULL REFERENCES equipments(id) ON DELETE CASCADE,
    PRIMARY KEY (rental_id, equipment_id)
);

-- migrate:down
DROP TABLE IF EXISTS rental_equipments;
DROP TABLE IF EXISTS rentals;
