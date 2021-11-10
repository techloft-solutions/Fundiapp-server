CREATE TABLE provider_professions(
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    provider_id VARCHAR(255) NOT NULL,
    category_id BIGINT(20) NOT NULL
);  