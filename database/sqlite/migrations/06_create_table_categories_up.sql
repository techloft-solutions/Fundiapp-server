CREATE TABLE categories(
    id INT(20) PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    parent_id INT(20),
    description TEXT,
    level INT(20) DEFAULT 0,
    icon_url VARCHAR(255) NOT NULL,
    industry_id INTEGER,
    FOREIGN KEY (parent_id) references categories(id),
    FOREIGN KEY (industry_id) references industries(id)
);