CREATE TABLE user_locations(
    id BIGINT(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(255) NOT NULL,
    location_id VARCHAR(255) NOT NULL
);