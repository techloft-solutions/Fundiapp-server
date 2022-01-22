CREATE TABLE IF NOT EXISTS subscriptions(
    id INT(20) PRIMARY KEY AUTO_INCREMENT,
    author_id VARCHAR(255) NOT NULL,
    comment TEXT NOT NULL,
    rating DECIMAL(1,1) NOT NULL,
    rating_integrity DECIMAL(1,1),
    plan_id VARCHAR(255) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
    -- FOREIGN KEY (provider_id) REFERENCES providers(provider_id) ON DELETE CASCADE,
    -- FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE CASCADE,
    -- UNIQUE (provider_id, service_id),
    -- FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);