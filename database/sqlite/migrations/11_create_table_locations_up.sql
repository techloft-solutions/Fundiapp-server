CREATE TABLE locations(
    location_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    title VARCHAR(255),
    city VARCHAR(255),
    state VARCHAR(255),
    country VARCHAR(255),
    latitude VARCHAR(255),
    longitude VARCHAR(255),
    booking_id VARCHAR(255)
);