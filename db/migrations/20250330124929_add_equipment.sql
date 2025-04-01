-- migrate:up
CREATE TABLE equipments (
    id VARCHAR(255) PRIMARY KEY,
    image_url TEXT,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    stock INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT,
    condition VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE camper_equipments (
    camper_id VARCHAR(255) NOT NULL REFERENCES campers(id) ON DELETE CASCADE,
    equipment_id VARCHAR(255) NOT NULL REFERENCES equipments(id) ON DELETE CASCADE,
    PRIMARY KEY (camper_id, equipment_id)
);

-- migrate:down
DROP TABLE IF EXISTS camper_equipments;
DROP TABLE IF EXISTS equipments;
