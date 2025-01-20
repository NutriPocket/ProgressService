CREATE DATABASE IF NOT EXISTS mydb; 

USE mydb;

SET time_zone = '+00:00';

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE TABLE IF NOT EXISTS jwt_blacklist (
    signature VARCHAR(100) PRIMARY KEY,
    expires_at TIMESTAMP NOT NULL,
    INDEX idx_expires_at (expires_at)
);

-- Delete expired JWTs

DELIMITER //

CREATE EVENT IF NOT EXISTS delete_expired_blacklist
ON SCHEDULE EVERY 1 MINUTE
DO
BEGIN
    DELETE FROM jwt_blacklist WHERE expires_at < NOW();
END //

DELIMITER ;