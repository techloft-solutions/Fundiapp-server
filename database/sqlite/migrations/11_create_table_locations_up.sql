CREATE TABLE locations(
    location_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    name VARCHAR(255),
    city VARCHAR(255),
    state VARCHAR(255),
    country VARCHAR(255),
    latitude VARCHAR(255) NOT NULL,
    longitude VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);