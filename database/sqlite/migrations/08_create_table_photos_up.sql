CREATE TABLE photos(
    photo_id VARCHAR(255) PRIMARY KEY,
    uploaded_by VARCHAR(255) NOT NULL,
    photo_url VARCHAR(255) NOT NULL,
    booking_id VARCHAR(255),
    portfolio_id VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (uploaded_by) REFERENCES users(user_id)
);