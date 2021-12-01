CREATE TABLE portfolios(
    portfolio_id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    owner_id VARCHAR(255) NOT NULL,
    booking_id VARCHAR(255) DEFAULT NULL,
    service_id INT(20) DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (booking_id) REFERENCES bookings(booking_id),
    FOREIGN KEY (service_id) REFERENCES services(id),
    FOREIGN KEY (owner_id) REFERENCES providers(provider_id)
);