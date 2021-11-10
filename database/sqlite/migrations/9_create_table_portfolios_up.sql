CREATE TABLE portfolios(
    portfolio_id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    booking_id VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);