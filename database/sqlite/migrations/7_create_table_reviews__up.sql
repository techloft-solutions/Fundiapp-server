CREATE TABLE IF NOT EXISTS reviews(
    review_id VARCHAR(255) PRIMARY KEY,
    author_id VARCHAR(255) NOT NULL,
    comment TEXT NOT NULL,
    rating DECIMAL(1,1) NOT NULL,
    rating_quality DECIMAL(1,1),
    rating_resposiveness DECIMAL(1,1),
    rating_integrity DECIMAL(1,1),
    rating_competence DECIMAL(1,1),
    provider_id VARCHAR(255) NOT NULL,
    service_id VARCHAR(255) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
    -- FOREIGN KEY (provider_id) REFERENCES providers(provider_id) ON DELETE CASCADE,
    -- FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE CASCADE,
    -- UNIQUE (provider_id, service_id),
    -- FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);