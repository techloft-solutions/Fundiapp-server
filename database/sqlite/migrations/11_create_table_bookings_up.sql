CREATE TABLE IF NOT EXISTS bookings (
    booking_id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    description TEXT(255),
    client_id VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255),
    location_id VARCHAR(255) NOT NULL,
    service_id VARCHAR(255),
    start_date DATETIME NOT NULL,
    status VARCHAR(255) NOT NULL,
    urgent BOOLEAN DEFAULT FALSE,
    is_request BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (provider_id) REFERENCES providers(provider_id),
    FOREIGN KEY (location_id) REFERENCES locations(location_id)
);