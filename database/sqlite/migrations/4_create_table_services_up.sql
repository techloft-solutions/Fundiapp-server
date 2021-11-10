CREATE TABLE services(
    service_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER,
    provider_id VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);