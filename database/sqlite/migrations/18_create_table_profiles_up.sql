/*CREATE TABLE IF NOT EXISTS profiles (
    profile_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    photo_url VARCHAR(255),
    location_id VARCHAR(255),
    account_type VARCHAR(255) DEFAULT 'client',
    status VARCHAR(255) DEFAULT 'active',
    verified BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
);*/
