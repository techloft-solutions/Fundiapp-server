CREATE TABLE providers(
    provider_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    location_id VARCHAR(255) NOT NULL,
    ratings_average INT(11) DEFAULT 0,
    reviews_count INT(11) DEFAULT 0,
    services_count INT(11) DEFAULT 0,
    jobs_count INT(11) DEFAULT 0,
    portfolio_count INT(11) DEFAULT 0,
    rate_per_hour VARCHAR(255),
    rate_per_unit VARCHAR(255),
    currency VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
);