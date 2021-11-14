CREATE TABLE categories(
    id INT(20) PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    parent_id INT(20),
    profession VARCHAR(255),
    description TEXT
);