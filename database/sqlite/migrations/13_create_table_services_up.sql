CREATE TABLE services(
    id INT(20) PRIMARY KEY AUTO_INCREMENT,
    -- user_id VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price INTEGER,
    description TEXT,
    currency CHAR(3) DEFAULT 'KES',
    price_unit VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (provider_id) REFERENCES providers(provider_id)
);