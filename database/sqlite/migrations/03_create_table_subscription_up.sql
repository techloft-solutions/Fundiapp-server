CREATE TABLE IF NOT EXISTS subscriptions (
    subscription_id VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL,
    payment_id VARCHAR(255) NOT NULL,
    plan_id VARCHAR(255) NOT NULL,
    auto_renew BOOLEAN DEFAULT FALSE,
    status VARCHAR(255),
    billing_cycles INT NULL,
    next_billing_at DATETIME,
    activated_at DATETIME,
    cancelled_at DATETIME,
    starts_at DATETIME DEFAULT NULL,
    expires_at DATETIME DEFAULT NULL,
    ends_at DATETIME DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
    -- FOREIGN KEY (provider_id) REFERENCES providers(provider_id) ON DELETE CASCADE,
    -- FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE CASCADE,
    -- UNIQUE (provider_id, service_id),
    -- FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);