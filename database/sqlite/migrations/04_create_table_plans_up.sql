CREATE TABLE IF NOT EXISTS plans(
    id INT(20) PRIMARY KEY AUTO_INCREMENT
    -- FOREIGN KEY (provider_id) REFERENCES providers(provider_id) ON DELETE CASCADE,
    -- FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE CASCADE,
    -- UNIQUE (provider_id, service_id),
    -- FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);